package responder

import (
	"context"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/mocks"
	responder_mocks "github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder/mocks"
	"github.com/stretchr/testify/assert"
	"net/netip"
	"testing"
	"time"
)

type testResponderMocks struct {
	peerID    protocol.ID
	table     *ktable_mocks.TableBatch
	responder responder
	sender    dht.NodeInfo
	limiter   *responder_mocks.Limiter
}

func newTestResponderMocks(t *testing.T) testResponderMocks {
	peerID := protocol.RandomNodeID()
	tableMock := ktable_mocks.NewTableBatch(t)
	limiter := responder_mocks.NewLimiter(t)
	return testResponderMocks{
		peerID: peerID,
		table:  tableMock,
		responder: responder{
			peerID:                   peerID,
			kTable:                   tableMock,
			sampleInfoHashesInterval: 20,
			limiter:                  limiter,
		},
		sender:  dht.RandomNodeInfo(4),
		limiter: limiter,
	}
}

var _ ktable.Peer = mockedPeer{}

type mockedPeer struct {
	dht.NodeInfo
}

func (m mockedPeer) ID() protocol.ID {
	return m.NodeInfo.ID
}

func (m mockedPeer) Addr() netip.AddrPort {
	return m.NodeInfo.Addr.ToAddrPort()
}

func (m mockedPeer) Time() time.Time {
	return time.Time{}
}

func (m mockedPeer) Dropped() bool {
	return false
}

func (m mockedPeer) IsSampleInfoHashesCandidate() bool {
	return true
}

type mockedHash struct {
	id        protocol.ID
	nodeInfos []dht.NodeInfo
}

func (m mockedHash) ID() protocol.ID {
	return m.id
}

func (m mockedHash) Peers() []ktable.HashPeer {
	var peers []ktable.HashPeer
	for _, nodeInfo := range m.nodeInfos {
		peers = append(peers, ktable.HashPeer{
			//ID:   nodeInfo.ID,
			Addr: nodeInfo.Addr.ToAddrPort(),
		})
	}
	return peers
}

func (m mockedHash) Dropped() bool {
	return false
}

func expectLimiterAllow(m testResponderMocks) {
	m.limiter.On("Allow", m.sender.Addr.ToAddrPort().Addr()).Return(true)
}

//func expectPutPeer(m testResponderMocks) {
//	m.table.On("BatchCommand", ktable.PutPeer{ID: m.sender.ID, Addr: m.sender.Addr.ToAddrPort()}).Return(nil)
//}

func TestResponder_ping(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "ping",
			A: &dht.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{ID: mocks.peerID}, ret)
	assert.NoError(t, err)
}

func TestResponder_ping__missing_args(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "ping",
		},
	}
	expectLimiterAllow(mocks)
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_find_node(t *testing.T) {
	mocks := newTestResponderMocks(t)
	target := protocol.RandomNodeID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "find_node",
			A: &dht.MsgArgs{
				ID:     mocks.sender.ID,
				Target: target,
			},
		},
	}
	nodes := dht.CompactIPv4NodeInfo{
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
	}
	peers := []ktable.Peer{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
		mockedPeer{nodes[2]},
	}
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	mocks.table.On("GetClosestPeers", target).Return(peers)
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID:    mocks.peerID,
		Nodes: nodes,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_find_node__missing_target(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "find_node",
			A: &dht.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	expectLimiterAllow(mocks)
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_get_peers__values(t *testing.T) {
	mocks := newTestResponderMocks(t)
	infoHash := protocol.RandomNodeID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "get_peers",
			A: &dht.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
			},
		},
	}
	nodeInfos := []dht.NodeInfo{
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	mocks.table.On("GetHashOrClosestPeers", infoHash).Return(ktable.GetHashOrClosestPeersResult{
		Hash:  mockedHash{nodeInfos: nodeInfos},
		Found: true,
	})
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.peerID,
		Values: []dht.NodeAddr{
			nodeInfos[0].Addr,
			nodeInfos[1].Addr,
			nodeInfos[2].Addr,
		},
		Nodes: nil,
		Token: &expectedToken,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_get_peers__nodes(t *testing.T) {
	mocks := newTestResponderMocks(t)
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	infoHash := protocol.RandomNodeID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "get_peers",
			A: &dht.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
			},
		},
	}
	nodes := dht.CompactIPv4NodeInfo{
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
	}
	peers := []ktable.Peer{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
		mockedPeer{nodes[2]},
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())
	mocks.table.On("GetHashOrClosestPeers", infoHash).Return(ktable.GetHashOrClosestPeersResult{
		ClosestPeers: peers,
	})
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID:     mocks.peerID,
		Values: nil,
		Nodes:  nodes,
		Token:  &expectedToken,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_get_peers__missing_info_hash(t *testing.T) {
	mocks := newTestResponderMocks(t)
	expectLimiterAllow(mocks)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "get_peers",
			A: &dht.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_announce_peer__implied_port(t *testing.T) {
	mocks := newTestResponderMocks(t)
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	infoHash := protocol.RandomNodeID()
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "announce_peer",
			A: &dht.MsgArgs{
				ID:          mocks.sender.ID,
				InfoHash:    infoHash,
				Token:       expectedToken,
				ImpliedPort: true,
			},
		},
	}
	mocks.table.On("BatchCommand", ktable.PutHash{
		ID: infoHash,
		Peers: []ktable.HashPeer{{
			//ID:   mocks.sender.ID,
			Addr: mocks.sender.Addr.ToAddrPort(),
		}},
	}).Return(nil)
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.peerID,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_announce_peer__specified_port(t *testing.T) {
	mocks := newTestResponderMocks(t)
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	infoHash := protocol.RandomNodeID()
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())
	port := krpc.RandomNodeInfo(4).Addr.Port
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "announce_peer",
			A: &dht.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
				Token:    expectedToken,
				Port:     &port,
			},
		},
	}
	mocks.table.On("BatchCommand", ktable.PutHash{
		ID: infoHash,
		Peers: []ktable.HashPeer{{
			//ID:   mocks.sender.ID,
			Addr: netip.AddrPortFrom(mocks.sender.Addr.ToAddrPort().Addr(), uint16(port)),
		}},
	}).Return(nil)
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.peerID,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_sample_infohashes(t *testing.T) {
	mocks := newTestResponderMocks(t)
	expectLimiterAllow(mocks)
	//expectPutPeer(mocks)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "sample_infohashes",
			A: &dht.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	infoHashes := []ktable.Hash{
		mockedHash{id: protocol.RandomNodeID()},
		mockedHash{id: protocol.RandomNodeID()},
		mockedHash{id: protocol.RandomNodeID()},
	}
	nodes := dht.CompactIPv4NodeInfo{
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
	}
	peers := []ktable.Peer{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
		mockedPeer{nodes[2]},
	}
	num := int64(123)
	mocks.table.On("SampleHashesAndPeers").Return(ktable.SampleHashesAndPeersResult{
		Hashes:      infoHashes,
		Peers:       peers,
		TotalHashes: int(num),
	})
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID:    mocks.peerID,
		Nodes: nodes,
		Bep51Return: dht.Bep51Return{
			Samples: &dht.CompactInfohashes{
				infoHashes[0].ID(),
				infoHashes[1].ID(),
				infoHashes[2].ID(),
			},
			Num:      &num,
			Interval: &mocks.responder.sampleInfoHashesInterval,
		},
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_unknown_method(t *testing.T) {
	mocks := newTestResponderMocks(t)
	expectLimiterAllow(mocks)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "foo",
			A: &dht.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, err, ErrMethodUnknown)
}
