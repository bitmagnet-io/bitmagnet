package responder

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

type Responder interface {
	Respond(context.Context, dht.RecvMsg) (dht.Return, error)
}

type responder struct {
	nodeID                   protocol.ID
	kTable                   ktable.Table
	kTable6                  ktable.Table
	tokenSecret              []byte
	sampleInfoHashesInterval int64
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
	if args == nil {
		err = ErrMissingArguments
		return
	}

	ret.ID = r.nodeID

	switch msg.Msg.Q {
	case dht.QPing:
	case dht.QFindNode:
		if args.Target == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		wantsV4, wantsV6 := r.wants(args, msg.From)
		if wantsV4 {
			ret.Nodes = nodeInfosFromNodes(r.kTable.GetClosestNodes(args.Target)...)
		}
		if wantsV6 {
			ret.Nodes6 = nodeInfosFromNodes(r.kTable6.GetClosestNodes(args.Target)...)
		}
	case dht.QGetPeers:
		if args.InfoHash == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		wantsV4, wantsV6 := r.wants(args, msg.From)

		// Query v4 table if needed (for v4 nodes or values from a v4 request)
		if wantsV4 || msg.From.Addr().Is4() {
			result := r.kTable.GetHashOrClosestNodes(args.InfoHash)
			if msg.From.Addr().Is4() && result.Found {
				ret.Values = peersToNodeAddrs(result.Hash.Peers())
			}
			if wantsV4 {
				ret.Nodes = nodeInfosFromNodes(result.ClosestNodes...)
			}
		}

		// Query v6 table if needed (for v6 nodes or values from a v6 request)
		if wantsV6 || msg.From.Addr().Is6() {
			result := r.kTable6.GetHashOrClosestNodes(args.InfoHash)
			if msg.From.Addr().Is6() && result.Found {
				ret.Values = peersToNodeAddrs(result.Hash.Peers())
			}
			if wantsV6 {
				ret.Nodes6 = nodeInfosFromNodes(result.ClosestNodes...)
			}
		}

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

		r.tableForAddr(msg.From.Addr()).BatchCommand(ktable.PutHash{ID: args.InfoHash, Peers: []ktable.HashPeer{{
			Addr: netip.AddrPortFrom(msg.From.Addr(), msg.AnnouncePort()),
		}}})
	case dht.QSampleInfohashes:
		result := r.tableForAddr(msg.From.Addr()).SampleHashesAndNodes()
		samples := make(dht.CompactInfohashes, 0, len(result.Hashes))

		for _, h := range result.Hashes {
			samples = append(samples, h.ID())
		}

		ret.Samples = &samples
		if msg.From.Addr().Is4() {
			ret.Nodes = nodeInfosFromNodes(result.Nodes...)
		} else if msg.From.Addr().Is6() {
			ret.Nodes6 = nodeInfosFromNodes(result.Nodes...)
		}
		numInt64 := int64(result.TotalHashes)
		ret.Num = &numInt64
		ret.Interval = &r.sampleInfoHashesInterval
	default:
		err = ErrMethodUnknown
		return
	}

	return
}

func (r responder) tableForAddr(addr netip.Addr) ktable.Table {
	if addr.Is6() {
		return r.kTable6
	}
	return r.kTable
}

func (r responder) wants(args *dht.MsgArgs, from netip.AddrPort) (v4, v6 bool) {
	if len(args.Want) > 0 {
		for _, w := range args.Want {
			switch w {
			case dht.WantNodes:
				v4 = true
			case dht.WantNodes6:
				v6 = true
			}
		}
	}
	v4 = v4 || from.Addr().Is4()
	v6 = v6 || from.Addr().Is6()

	return
}

func peersToNodeAddrs(peers []ktable.HashPeer) []dht.NodeAddr {
	values := make([]dht.NodeAddr, 0, len(peers))
	for _, p := range peers {
		values = append(values, dht.NewNodeAddrFromAddrPort(p.Addr))
	}
	return values
}

// announceToken returns the token for an announce_peer request.
// A "token" key is included in the get_peers return value.
// The token value is a required argument for a future announce_peer query.
// The token value should be a short binary string.
// The queried node must verify that the token was previously sent to the same IP address as the querying node.
// Then the queried node should store the IP address of the querying node and the supplied port number
// under the infohash in its store of peer contact information.
// https://www.bittorrent.org/beps/bep_0005.html
func (r responder) announceToken(infoHash protocol.ID, nodeID protocol.ID, nodeAddr netip.Addr) string {
	bytes := r.tokenSecret
	bytes = append(bytes, r.nodeID[:]...)
	bytes = append(bytes, infoHash[:]...)
	bytes = append(bytes, nodeID[:]...)
	bytes = append(bytes, []byte(nodeAddr.String())...)
	tokenHash := md5.Sum(bytes)

	return hex.EncodeToString(tokenHash[:])
}

func nodeInfosFromNodes(ns ...ktable.Node) []dht.NodeInfo {
	if len(ns) == 0 {
		return nil
	}

	nodes := make([]dht.NodeInfo, 0, len(ns))
	for _, n := range ns {
		nodes = append(nodes, nodeInfoFromNode(n))
	}

	return nodes
}

func nodeInfoFromNode(n ktable.Node) dht.NodeInfo {
	return dht.NodeInfo{
		ID:   n.ID(),
		Addr: dht.NewNodeAddrFromAddrPort(n.Addr()),
	}
}
