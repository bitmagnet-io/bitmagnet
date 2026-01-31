package metainforequester

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"net/netip"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/peer_protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
)

type Requester interface {
	Request(context.Context, protocol.ID, netip.AddrPort) (Response, error)
}

type requester struct {
	clientID       protocol.ID
	dialTimeout    *atomic.Value[DialTimeout]
	requestTimeout *atomic.Value[RequestTimeout]
}

type ExtensionBit uint

// https://www.bittorrent.org/beps/bep_0004.html
// https://wiki.theory.org/BitTorrentSpecification.html#Reserved_Bytes
const (
	// ExtensionBitDht http://www.bittorrent.org/beps/bep_0005.html
	ExtensionBitDht = 0
	// ExtensionBitFast http://www.bittorrent.org/beps/bep_0006.html
	ExtensionBitFast = 2
	// ExtensionBitV2 "Hybrid torrent legacy to v2 upgrade"
	ExtensionBitV2                           = 7
	ExtensionBitAzureusExtensionNegotiation1 = 16
	ExtensionBitAzureusExtensionNegotiation2 = 17
	// ExtensionBitLtep LibTorrent Extension Protocol, http://www.bittorrent.org/beps/bep_0010.html
	ExtensionBitLtep = 20
	// ExtensionBitLocationAwareProtocol https://wiki.theory.org/BitTorrent_Location-aware_Protocol_1
	ExtensionBitLocationAwareProtocol = 43
	// ExtensionBitAzureusMessagingProtocol https://www.bittorrent.org/beps/bep_0004.html
	ExtensionBitAzureusMessagingProtocol = 63
)

type PeerExtensionBits [8]byte

func NewPeerExtensionBits(bits ...ExtensionBit) (ret PeerExtensionBits) {
	for _, b := range bits {
		ret = ret.WithBit(b, true)
	}

	return
}

func (pex PeerExtensionBits) WithBit(bit ExtensionBit, on bool) PeerExtensionBits {
	if on {
		pex[7-bit/8] |= 1 << (bit % 8)
	} else {
		pex[7-bit/8] &^= 1 << (bit % 8)
	}

	return pex
}

func (pex PeerExtensionBits) GetBit(bit ExtensionBit) bool {
	return pex[7-bit/8]&(1<<(bit%8)) != 0
}

type HandshakeInfo struct {
	PeerID            protocol.ID
	PeerExtensionBits PeerExtensionBits
}

type Response struct {
	HandshakeInfo
	Info metainfo.Info
}

func (r *requester) Request(ctx context.Context, infoHash protocol.ID, addr netip.AddrPort) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.requestTimeout.Get()))
	defer cancel()

	conn, err := r.connect(ctx, addr)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrConnect, err)
	}

	defer func() {
		_ = conn.Close()
	}()

	hsInfo, err := btHandshake(conn, infoHash, r.clientID)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrHandshake, err)
	}

	metadataSize, utMetadata, err := exHandshake(conn)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrExHandshake, err)
	}

	if err = requestAllPieces(conn, metadataSize, utMetadata); err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrRequestPieces, err)
	}

	pieces, err := readAllPieces(conn, metadataSize)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrReadPieces, err)
	}

	parsed, err := metainfo.ParseMetaInfoBytes(infoHash, pieces)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrParse, err)
	}

	return Response{
		HandshakeInfo: hsInfo,
		Info:          parsed,
	}, nil
}

func (r requester) connect(ctx context.Context, addr netip.AddrPort) (*net.TCPConn, error) {
	dialer := &net.Dialer{
		Timeout:   time.Duration(r.dialTimeout.Get()),
		KeepAlive: -1,
	}

	c, err := dialer.DialContext(ctx, "tcp4", addr.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDial, err)
	}

	tcpConn := c.(*net.TCPConn)
	closeConn := func() {
		_ = tcpConn.Close()
	}
	// If sec == 0, operating system discards any unsent or unacknowledged data [after Close() has been called].
	if err := tcpConn.SetLinger(0); err != nil {
		closeConn()

		return nil, fmt.Errorf("%w: %w", ErrSetLinger, err)
	}

	deadline, ok := ctx.Deadline()
	if ok {
		if err := tcpConn.SetDeadline(deadline); err != nil {
			closeConn()

			return nil, fmt.Errorf("%w: %w", ErrSetDeadline, err)
		}
	}

	return tcpConn, nil
}

var myExBits = NewPeerExtensionBits(ExtensionBitDht, ExtensionBitLtep)

func btHandshake(rw io.ReadWriter, infoHash protocol.ID, clientID protocol.ID) (HandshakeInfo, error) {
	handshakeBytes := make([]byte, 0, 68)
	handshakeBytes = append(handshakeBytes, peer_protocol.Protocol...)
	handshakeBytes = append(handshakeBytes, myExBits[:]...)
	handshakeBytes = append(handshakeBytes, infoHash[:]...)
	handshakeBytes = append(handshakeBytes, clientID[:]...)

	if _, err := rw.Write(handshakeBytes); err != nil {
		return HandshakeInfo{}, fmt.Errorf("%w: %w", ErrWrite, err)
	}

	handshakeResponse := make([]byte, 68)
	if _, err := io.ReadFull(rw, handshakeResponse); err != nil {
		return HandshakeInfo{}, fmt.Errorf("%w: %w", ErrRead, err)
	}

	if !bytes.HasPrefix(handshakeResponse, []byte(peer_protocol.Protocol)) {
		return HandshakeInfo{}, ErrInvalidResponse
	}

	var peerExBits PeerExtensionBits

	copy(peerExBits[:], handshakeResponse[20:28])

	if !peerExBits.GetBit(ExtensionBitLtep) {
		return HandshakeInfo{}, ErrUnsupported
	}

	var resHash protocol.ID

	copy(resHash[:], handshakeResponse[28:48])

	if resHash != infoHash {
		return HandshakeInfo{}, ErrInfoHashMismatch
	}

	var resPeerID protocol.ID

	copy(resPeerID[:], handshakeResponse[48:68])

	return HandshakeInfo{
		PeerID:            resPeerID,
		PeerExtensionBits: peerExBits,
	}, nil
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

func exHandshake(rw io.ReadWriter) (uint, uint8, error) {
	if _, err := rw.Write([]byte("\x00\x00\x00\x1a\x14\x00d1:md11:ut_metadatai1eee")); err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrWrite, err)
	}

	rExMessage, err := readExMessage(rw)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrRead, err)
	}

	// Extension Handshake has the Extension Message id = 0x00
	if rExMessage[1] != 0 {
		return 0, 0, fmt.Errorf("%w: %w", ErrInvalidResponse, ErrFirstExMessage)
	}

	rRootDict := new(rootDict)
	if err := bencode.Unmarshal(rExMessage[2:], rRootDict); err != nil {
		return 0, 0, fmt.Errorf("%w: %w: %w", ErrInvalidResponse, ErrUnmarshal, err)
	}

	if 0 >= rRootDict.MetadataSize || rRootDict.MetadataSize >= maxMetadataSize {
		return 0, 0, fmt.Errorf("%w: %w", ErrInvalidResponse, ErrSize)
	}

	if 0 >= rRootDict.M.UTMetadata || rRootDict.M.UTMetadata >= 255 {
		return 0, 0, fmt.Errorf("%w: %w", ErrInvalidResponse, ErrUTMetadata)
	}

	return uint(rRootDict.MetadataSize), uint8(rRootDict.M.UTMetadata), nil
}

func requestAllPieces(w io.Writer, metadataSize uint, utMetadata uint8) error {
	nPieces := int(math.Ceil(float64(metadataSize) / math.Pow(2, 14)))
	for piece := range nPieces {
		extDictDump, err := bencode.Marshal(extDict{
			MsgType: 0,
			Piece:   piece,
		})
		if err != nil { // ASSERT
			panic(err)
		}

		if _, err = fmt.Fprintf(w,
			"%s\x14%s%s",
			uintToBigEndian4(uint(2+len(extDictDump))),
			[]byte{utMetadata},
			extDictDump); err != nil {
			return err
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
		// run TestDecoder() function in leech_test.go in case you have any doubts.
		rMessageBuf := bytes.NewBuffer(rUmMessage[2:])
		rExtDict := new(extDict)

		if err = bencode.NewDecoder(rMessageBuf).Decode(rExtDict); err != nil {
			return nil, fmt.Errorf("%w: %w: %w", ErrInvalidResponse, ErrUnmarshal, err)
		}

		if rExtDict.MsgType == 2 { // reject
			return nil, ErrRejected
		}

		if rExtDict.MsgType == 1 { // data
			// Get the unread bytes!
			metadataPiece := rMessageBuf.Bytes()
			// BEP 9 explicitly states:
			//   > If the piece is the last piece of the metadata, it may be less than 16kiB. If
			//   > it is not the last piece of the metadata, it MUST be 16kiB.
			//
			// Hence...
			//   ... if the length of metadataPiece is more than 16kiB, we err.
			if len(metadataPiece) > 16*1024 {
				return nil, fmt.Errorf("%w: %w: > 16kiB", ErrInvalidResponse, ErrSize)
			}

			receivedSize += uint(len(metadataPiece))
			// ... if the length of @metadataPiece is less than 16kiB AND metadataBytes is NOT
			// complete then we err.
			if len(metadataPiece) < 16*1024 && receivedSize != metadataSize {
				return nil, fmt.Errorf("%w: %w: < 16 kiB but incomplete", ErrInvalidResponse, ErrSize)
			}

			if receivedSize > metadataSize {
				return nil, fmt.Errorf("%w: %w: received > handshake", ErrInvalidResponse, ErrSize)
			}

			piece := rExtDict.Piece
			copy(
				metadataBytes[piece*int(math.Pow(2, 14)):piece*int(math.Pow(2, 14))+len(metadataPiece)],
				metadataPiece,
			)
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
		return nil, fmt.Errorf("%w: %w: > %d", ErrInvalidResponse, ErrSize, maxMetadataSize)
	}

	messageBytes := make([]byte, length)
	if _, err := io.ReadFull(r, messageBytes); err != nil {
		return nil, err
	}

	return messageBytes, nil
}
