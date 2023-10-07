package responder

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/anacrolix/dht/v2/transactions"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/routingtable"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	PeerID krpc.ID `name:"dht_peer_id"`
	Table  routingtable.Table
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Responder Responder
}

func New(p Params) Result {
	return Result{
		Responder: responder{
			peerID:                   p.PeerID,
			table:                    p.Table,
			tokenSecret:              []byte(transactions.DefaultIdIssuer.Issue()),
			sampleInfoHashesInterval: 10,
			logger:                   p.Logger.Named("dht_responder"),
		},
	}
}

type Responder interface {
	Respond(context.Context, dht.RecvMsg) (krpc.Return, error)
}

type responder struct {
	peerID                   [20]byte
	table                    routingtable.Table
	tokenSecret              []byte
	sampleInfoHashesInterval int64
	logger                   *zap.SugaredLogger
}

var ErrMissingArguments = krpc.Error{
	Code: krpc.ErrorCodeProtocolError,
	Msg:  "missing arguments",
}

var ErrInvalidToken = krpc.Error{
	Code: krpc.ErrorCodeProtocolError,
	Msg:  "invalid token",
}

var ErrMethodUnknown = krpc.Error{
	Code: krpc.ErrorCodeMethodUnknown,
	Msg:  "method Unknown",
}

func (r responder) Respond(_ context.Context, msg dht.RecvMsg) (ret krpc.Return, err error) {
	args := msg.Msg.A
	defer (func() {
		if err != nil {
			r.logger.Debugw("responding with error", "q", msg.Msg.Q, "args", args, "error", err)
		}
	})()
	if args == nil {
		err = ErrMissingArguments
		return
	}
	ret.ID = r.peerID
	nodeInfo := krpc.NodeInfo{
		ID:   args.ID,
		Addr: msg.From,
	}
	switch msg.Msg.Q {
	case dht.QPing:
		r.logger.Debug(dht.QPing)
	case dht.QFindNode:
		if args.Target == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		ret.Nodes = r.table.FindNode(args.Target)
		r.logger.Debugw(dht.QFindNode, "n", len(ret.Nodes))
	case dht.QGetPeers:
		if args.InfoHash == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		values, nodes := r.table.GetPeers(args.InfoHash)
		token := r.announceToken(args.InfoHash, nodeInfo)
		ret.Token = &token
		ret.Values = values
		ret.Nodes = nodes
		r.logger.Debugw(dht.QGetPeers, "nValues", len(ret.Values), "nNodes", len(ret.Nodes))
	case dht.QAnnouncePeer:
		if args.InfoHash == [20]byte{} {
			err = ErrMissingArguments
			return
		}
		if args.Token != r.announceToken(args.InfoHash, nodeInfo) {
			err = ErrInvalidToken
			return
		}
		r.table.ReceivePeersForHash(args.InfoHash, krpc.NodeAddr{
			IP:   nodeInfo.Addr.IP,
			Port: msg.AnnouncePort(),
		})
		r.logger.Debug(dht.QAnnouncePeer)
	case dht.QSampleInfohashes:
		hashes, addrs, num := r.table.SampleInfoHashes()
		ret.Samples = &hashes
		ret.Nodes = addrs
		ret.Num = &num
		ret.Interval = &r.sampleInfoHashesInterval
		r.logger.Debugw(dht.QSampleInfohashes, "nSamples", len(hashes), "nNodes", len(ret.Nodes), "nTotalHashes", num)
	default:
		err = ErrMethodUnknown
		return
	}
	r.table.ReceiveNodeInfo(nodeInfo)
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
func (r responder) announceToken(infoHash krpc.ID, nodeInfo krpc.NodeInfo) string {
	bytes := r.tokenSecret
	bytes = append(bytes, r.peerID[:]...)
	bytes = append(bytes, infoHash[:]...)
	bytes = append(bytes, nodeInfo.ID[:]...)
	bytes = append(bytes, []byte(nodeInfo.Addr.IP.String())...)
	tokenHash := md5.Sum(bytes)
	return hex.EncodeToString(tokenHash[:])
}
