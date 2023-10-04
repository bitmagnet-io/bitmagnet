package metainforequester

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/peer_protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"io"
	"math"
	"net"
	"time"
)

type Params struct {
	fx.In
	Config Config
}

type Result struct {
	fx.Out
	Requester Requester
}

func New(p Params) Result {
	return Result{
		Requester: requester{
			clientID: dht.RandomPeerID(),
			timeout:  p.Config.RequestTimeout,
			dialer: &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: -1,
			},
			semaphore: semaphore.NewWeighted(int64(p.Config.MaxConcurrency)),
		},
	}
}

type Requester interface {
	Request(ctx context.Context, infoHash model.Hash20, node krpc.NodeAddr) (metainfo.Info, error)
}

type requester struct {
	clientID  model.Hash20
	timeout   time.Duration
	dialer    *net.Dialer
	semaphore *semaphore.Weighted
}

func (r requester) Request(ctx context.Context, infoHash model.Hash20, node krpc.NodeAddr) (metainfo.Info, error) {
	if semErr := r.semaphore.Acquire(ctx, 1); semErr != nil {
		return metainfo.Info{}, semErr
	}
	defer r.semaphore.Release(1)
	gotPieces := make(chan []byte)
	gotErr := make(chan error)
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	go (func() {
		conn, connErr := r.connect(timeoutCtx, node)
		if connErr != nil {
			gotErr <- connErr
			return
		}
		defer func() {
			_ = conn.Close()
		}()
		if btHandshakeErr := btHandshake(conn, infoHash, r.clientID); btHandshakeErr != nil {
			gotErr <- btHandshakeErr
			return
		}
		metadataSize, utMetadata, exHandshakeErr := exHandshake(conn)
		if exHandshakeErr != nil {
			gotErr <- exHandshakeErr
			return
		}
		if requestAllPiecesErr := requestAllPieces(conn, metadataSize, utMetadata); requestAllPiecesErr != nil {
			gotErr <- requestAllPiecesErr
			return
		}
		pieces, readAllPiecesErr := readAllPieces(conn, metadataSize)
		if readAllPiecesErr != nil {
			gotErr <- readAllPiecesErr
			return
		}
		gotPieces <- pieces
	})()
	select {
	case <-timeoutCtx.Done():
		return metainfo.Info{}, timeoutCtx.Err()
	case err := <-gotErr:
		return metainfo.Info{}, err
	case pieces := <-gotPieces:
		return metainfo.ParseMetaInfoBytes(infoHash, pieces)
	}
}

func (r requester) connect(ctx context.Context, addr krpc.NodeAddr) (conn *net.TCPConn, err error) {
	c, dialErr := r.dialer.DialContext(ctx, "tcp4", addr.String())
	if dialErr != nil {
		err = dialErr
		return
	}
	tcpConn := c.(*net.TCPConn)
	closeConn := func() {
		_ = tcpConn.Close()
	}
	// If sec == 0, operating system discards any unsent or unacknowledged data [after Close() has been called].
	if setLingerErr := tcpConn.SetLinger(0); setLingerErr != nil {
		err = setLingerErr
		closeConn()
		return
	}
	deadline, ok := ctx.Deadline()
	if ok {
		if setDeadlineErr := tcpConn.SetDeadline(deadline); setDeadlineErr != nil {
			err = setDeadlineErr
			closeConn()
			return
		}
	}
	return tcpConn, nil
}

func btHandshake(rw io.ReadWriter, infoHash model.Hash20, clientID model.Hash20) error {
	handshakeBytes := make([]byte, 0, 68)
	handshakeBytes = append(handshakeBytes, []byte(peer_protocol.Protocol+"\u0000\u0000\u0000\u0000\u0000\u0010\u0000\u0001")...)
	handshakeBytes = append(handshakeBytes, infoHash[:]...)
	handshakeBytes = append(handshakeBytes, clientID[:]...)
	if n, hsErr := rw.Write(handshakeBytes); hsErr != nil {
		return hsErr
	} else if n != 68 {
		panic("handshake bytes must have length 68")
	}
	handshakeResponse := make([]byte, 68)
	if n, hsErr := io.ReadFull(rw, handshakeResponse); hsErr != nil {
		return fmt.Errorf("failed to read all handshake bytes (%d): %w / %s", n, hsErr, infoHash.String())
	}
	if !bytes.HasPrefix(handshakeResponse, []byte(peer_protocol.Protocol)) {
		return errors.New("invalid handshake response received")
	}
	// TODO: maybe check for the infohash sent by the remote peer to double check?
	if (handshakeResponse[25] & 0x10) == 0 {
		return errors.New("peer does not support the extension protocol")
	}
	return nil
}

type rootDict struct {
	M            mDict `bencode:"m"`
	MetadataSize int   `bencode:"metadata_size"`
}

type mDict struct {
	UTMetadata int `bencode:"ut_metadata"`
}

type extDict struct {
	MsgType int `bencode:"msg_type"`
	Piece   int `bencode:"piece"`
}

const maxMetadataSize = 10 * 1024 * 1024

func exHandshake(rw io.ReadWriter) (metadataSize uint, utMetadata uint8, err error) {
	if _, writeErr := rw.Write([]byte("\x00\x00\x00\x1a\x14\x00d1:md11:ut_metadatai1eee")); err != nil {
		err = writeErr
		return
	}
	rExMessage, readErr := readExMessage(rw)
	if readErr != nil {
		err = readErr
		return
	}
	// Extension Handshake has the Extension Message ID = 0x00
	if rExMessage[1] != 0 {
		err = errors.New("first extension message is not an extension handshake")
		return
	}
	rRootDict := new(rootDict)
	if unmarshalErr := bencode.Unmarshal(rExMessage[2:], rRootDict); unmarshalErr != nil {
		err = unmarshalErr
		return
	}
	if !(0 < rRootDict.MetadataSize && rRootDict.MetadataSize < maxMetadataSize) {
		err = errors.New("metadata too big or its size is less than or equal zero")
		return
	}
	if !(0 < rRootDict.M.UTMetadata && rRootDict.M.UTMetadata < 255) {
		err = errors.New("ut_metadata is not an uint8")
		return
	}
	return uint(rRootDict.MetadataSize), uint8(rRootDict.M.UTMetadata), nil
}

func requestAllPieces(w io.Writer, metadataSize uint, utMetadata uint8) error {
	nPieces := int(math.Ceil(float64(metadataSize) / math.Pow(2, 14)))
	for piece := 0; piece < nPieces; piece++ {
		extDictDump, err := bencode.Marshal(extDict{
			MsgType: 0,
			Piece:   piece,
		})
		if err != nil { // ASSERT
			panic(err)
		}
		if _, writeErr := w.Write([]byte(fmt.Sprintf(
			"%s\x14%s%s",
			uintToBigEndian4(uint(2+len(extDictDump))),
			[]byte{utMetadata},
			extDictDump,
		))); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

func uintToBigEndian4(i uint) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

func readAllPieces(r io.Reader, metadataSize uint) ([]byte, error) {
	metadataBytes := make([]byte, metadataSize)
	receivedSize := uint(0)
	for receivedSize < metadataSize {
		rUmMessage, err := readUmMessage(r)
		if err != nil {
			return nil, err
		}
		// Run TestDecoder() function in leech_test.go in case you have any doubts.
		rMessageBuf := bytes.NewBuffer(rUmMessage[2:])
		rExtDict := new(extDict)
		if decodeErr := bencode.NewDecoder(rMessageBuf).Decode(rExtDict); decodeErr != nil {
			return nil, decodeErr
		}
		if rExtDict.MsgType == 2 { // reject
			return nil, errors.New("remote peer rejected sending metadataBytes")
		}
		if rExtDict.MsgType == 1 { // data
			// Get the unread bytes!
			metadataPiece := rMessageBuf.Bytes()
			// BEP 9 explicitly states:
			//   > If the piece is the last piece of the metadata, it may be less than 16kiB. If
			//   > it is not the last piece of the metadata, it MUST be 16kiB.
			//
			// Hence...
			//   ... if the length of @metadataPiece is more than 16kiB, we err.
			if len(metadataPiece) > 16*1024 {
				return nil, errors.New("metadataPiece > 16kiB")
			}
			piece := rExtDict.Piece
			copy(metadataBytes[piece*int(math.Pow(2, 14)):piece*int(math.Pow(2, 14))+len(metadataPiece)], metadataPiece)
			receivedSize += uint(len(metadataPiece))
			// ... if the length of @metadataPiece is less than 16kiB AND metadataBytes is NOT
			// complete then we err.
			if len(metadataPiece) < 16*1024 && receivedSize != metadataSize {
				return nil, errors.New("metadataPiece < 16 kiB but incomplete")
			}

			if receivedSize > metadataSize {
				return nil, errors.New("receivedSize > metadataSize")
			}
		}
	}
	return metadataBytes, nil
}

func readExMessage(r io.Reader) ([]byte, error) {
	for {
		rMessage, err := readMessage(r)
		if err != nil {
			return nil, err
		}
		// Every extension message has at least 2 bytes.
		if len(rMessage) < 2 {
			continue
		}
		// We are interested only in extension messages, whose first byte is always 20
		if rMessage[0] == byte(peer_protocol.Extended) {
			return rMessage, nil
		}
	}
}

// readUmMessage returns an ut_metadata extension message, sans the first 4 bytes indicating its
// length.
//
// It will IGNORE all non-"ut_metadata extension" messages!
func readUmMessage(r io.Reader) ([]byte, error) {
	for {
		rExMessage, err := readExMessage(r)
		if err != nil {
			return nil, err
		}

		if rExMessage[1] == 0x01 {
			return rExMessage, nil
		}
	}
}

// readMessage returns a BitTorrent message, sans the first 4 bytes indicating its length.
func readMessage(r io.Reader) ([]byte, error) {
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBytes); err != nil {
		return nil, err
	}
	length := uint(binary.BigEndian.Uint32(lengthBytes))

	// Some malicious/faulty peers say that they are sending a very long
	// message, and hence causing us to run out of memory.
	// This is a crude check that does not let it happen (i.e. boundary can probably be
	// tightened a lot more.)
	if length > maxMetadataSize {
		return nil, errors.New("message is longer than max allowed metadata size")
	}

	messageBytes := make([]byte, length)
	if _, err := io.ReadFull(r, messageBytes); err != nil {
		return nil, err
	}

	return messageBytes, nil
}
