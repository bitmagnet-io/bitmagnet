package responder

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"go.uber.org/fx"
	"golang.org/x/time/rate"
	"net/netip"
	"time"
)

type Params struct {
	fx.In
	PeerID protocol.ID `name:"peer_id"`
	KTable ktable.TableBatch
}

type Result struct {
	fx.Out
	Responder Responder
}

func New(p Params) Result {
	return Result{
		Responder: responder{
			peerID:                   p.PeerID,
			kTable:                   p.KTable,
			tokenSecret:              protocol.RandomNodeID().Bytes(),
			sampleInfoHashesInterval: 10,
			limiter:                  NewLimiter(rate.Every(time.Second/5), 20, rate.Every(time.Second), 10, 1000, time.Second*20),
		},
	}
}

type Responder interface {
	Respond(context.Context, dht.RecvMsg) (dht.Return, error)
}

type responder struct {
	peerID                   [20]byte
	kTable                   ktable.TableBatch
	tokenSecret              []byte
	sampleInfoHashesInterval int64
	limiter                  Limiter
}

var ErrMissingArguments = dht.Error{
	Code: dht.ErrorCodeProtocolError,
	Msg:  "missing arguments",
}

var ErrInvalidToken = dht.Error{
	Code: dht.ErrorCodeProtocolError,
	Msg:  "invalid token",
}

var ErrMethodUnknown = dht.Error{
	Code: dht.ErrorCodeMethodUnknown,
	Msg:  "method Unknown",
}

var ErrTooManyRequests = dht.Error{
	Code: dht.ErrorCodeGenericError,
	Msg:  "too many requests",
}

func (r responder) Respond(_ context.Context, msg dht.RecvMsg) (ret dht.Return, err error) {
	args := msg.Msg.A
	defer (func() {
		if err == nil {
			r.kTable.BatchCommand(ktable.PutPeer{ID: args.ID, Addr: msg.From})
		}
	})()
	// apply both overall and per-IP rate limiting:
	if !r.limiter.Allow(msg.From.Addr()) {
		err = ErrTooManyRequests
		return
	}
	if args == nil {
		err = ErrMissingArguments
		return
	}
	ret.ID = r.peerID
	switch msg.Msg.Q {
	case dht.QPing:
	case dht.QFindNode:
		if args.Target == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		closestPeers := r.kTable.GetClosestPeers(args.Target)
		ret.Nodes = nodeInfosFromPeers(closestPeers...)
	case dht.QGetPeers:
		if args.InfoHash == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		result := r.kTable.GetHashOrClosestPeers(args.InfoHash)
		if result.Found {
			hashPeers := result.Hash.Peers()
			values := make([]dht.NodeAddr, 0, len(hashPeers))
			for _, p := range hashPeers {
				values = append(values, dht.NewNodeAddrFromAddrPort(p.Addr))
			}
			ret.Values = values
		}
		ret.Nodes = nodeInfosFromPeers(result.ClosestPeers...)
		token := r.announceToken(args.InfoHash, args.ID, msg.From.Addr())
		ret.Token = &token
	case dht.QAnnouncePeer:
		if args.InfoHash == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		if args.Token != r.announceToken(args.InfoHash, args.ID, msg.From.Addr()) {
			err = ErrInvalidToken
			return
		}
		r.kTable.BatchCommand(ktable.PutHash{ID: args.InfoHash, Peers: []ktable.HashPeer{{
			ID:   args.ID,
			Addr: netip.AddrPortFrom(msg.From.Addr(), msg.AnnouncePort()),
		}}})
	case dht.QSampleInfohashes:
		result := r.kTable.SampleHashesAndPeers()
		samples := make(dht.CompactInfohashes, 0, len(result.Hashes))
		for _, h := range result.Hashes {
			samples = append(samples, h.ID())
		}
		ret.Samples = &samples
		ret.Nodes = nodeInfosFromPeers(result.Peers...)
		numInt64 := int64(result.TotalHashes)
		ret.Num = &numInt64
		ret.Interval = &r.sampleInfoHashesInterval
	default:
		err = ErrMethodUnknown
		return
	}
	return
}

// announceToken returns the token for an announce_peer request.
// A "token" key is included in the get_peers return value.
// The token value is a required argument for a future announce_peer query.
// The token value should be a short binary string.
// The queried node must verify that the token was previously sent to the same IP address as the querying node.
// Then the queried node should store the IP address of the querying node and the supplied port number under the infohash
// in its store of peer contact information.
// https://www.bittorrent.org/beps/bep_0005.html
func (r responder) announceToken(infoHash protocol.ID, nodeID protocol.ID, nodeAddr netip.Addr) string {
	bytes := r.tokenSecret
	bytes = append(bytes, r.peerID[:]...)
	bytes = append(bytes, infoHash[:]...)
	bytes = append(bytes, nodeID[:]...)
	bytes = append(bytes, []byte(nodeAddr.String())...)
	tokenHash := md5.Sum(bytes)
	return hex.EncodeToString(tokenHash[:])
}

func nodeInfosFromPeers(p ...ktable.Peer) []dht.NodeInfo {
	if len(p) == 0 {
		return nil
	}
	nodes := make([]dht.NodeInfo, 0, len(p))
	for _, peer := range p {
		nodes = append(nodes, nodeInfoFromPeer(peer))
	}
	return nodes
}

func nodeInfoFromPeer(p ktable.Peer) dht.NodeInfo {
	return dht.NodeInfo{
		ID:   p.ID(),
		Addr: dht.NewNodeAddrFromAddrPort(p.Addr()),
	}
}
