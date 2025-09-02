package responder

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	ktable_mocks "github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/mocks"
	"github.com/stretchr/testify/assert"
)

type testResponderMocks struct {
	nodeID    protocol.ID
	table     *ktable_mocks.Table
	table6    *ktable_mocks.Table
	responder responder
	sender    dht.NodeInfo
	sender6   dht.NodeInfo
}

func newTestResponderMocks(t *testing.T) testResponderMocks {
	t.Helper()

	nodeID := protocol.RandomNodeID()
	tableMock := ktable_mocks.NewTable(t)
	table6Mock := ktable_mocks.NewTable(t)

	return testResponderMocks{
		nodeID: nodeID,
		table:  tableMock,
		table6: table6Mock,
		responder: responder{
			nodeID:                   nodeID,
			kTable:                   tableMock,
			kTable6:                  table6Mock,
			tokenSecret:              protocol.RandomNodeID().Bytes(),
			sampleInfoHashesInterval: 20,
		},
		sender:  dht.RandomNodeInfo(4),
		sender6: dht.RandomNodeInfo(16),
	}
}

var _ ktable.Node = mockedPeer{}

type mockedPeer struct {
	dht.NodeInfo
}

func (m mockedPeer) ID() protocol.ID {
	return m.NodeInfo.ID
}

func (m mockedPeer) Addr() netip.AddrPort {
	return m.NodeInfo.Addr.ToAddrPort()
}

func (mockedPeer) Time() time.Time {
	return time.Time{}
}

func (mockedPeer) Dropped() bool {
	return false
}

func (mockedPeer) IsSampleInfoHashesCandidate() bool {
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
	peers := make([]ktable.HashPeer, 0, len(m.nodeInfos))
	for _, nodeInfo := range m.nodeInfos {
		peers = append(peers, ktable.HashPeer{
			// ID:   nodeInfo.ID,
			Addr: nodeInfo.Addr.ToAddrPort(),
		})
	}

	return peers
}

func (mockedHash) Dropped() bool {
	return false
}

func TestResponder_ping(t *testing.T) {
	t.Parallel()

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
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{ID: mocks.nodeID}, ret)
	assert.NoError(t, err)
}

func TestResponder_ping__missing_args(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "ping",
		},
	}
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_find_node(t *testing.T) {
	t.Parallel()

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
	peers := []ktable.Node{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
		mockedPeer{nodes[2]},
	}
	mocks.table.On("GetClosestNodes", target).Return(peers)
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID:    mocks.nodeID,
		Nodes: nodes,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_find_node__ipv6(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
	target := protocol.RandomNodeID()
	msg := dht.RecvMsg{
		From: mocks.sender6.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "find_node",
			A: &dht.MsgArgs{
				ID:     mocks.sender6.ID,
				Target: target,
			},
		},
	}
	nodes := dht.CompactIPv6NodeInfo{
		dht.RandomNodeInfo(16),
		dht.RandomNodeInfo(16),
	}
	peers := []ktable.Node{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
	}
	mocks.table6.On("GetClosestNodes", target).Return(peers)
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID:     mocks.nodeID,
		Nodes:  nil,
		Nodes6: nodes,
	}, ret)
	assert.NoError(t, err)
	mocks.table.AssertNotCalled(t, "GetClosestNodes")
}

func TestResponder_find_node__missing_target(t *testing.T) {
	t.Parallel()

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
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_get_peers__values(t *testing.T) {
	t.Parallel()

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
	closestNodes := []ktable.Node{
		mockedPeer{dht.RandomNodeInfo(4)},
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())
	mocks.table.On("GetHashOrClosestNodes", infoHash).Return(ktable.GetHashOrClosestNodesResult{
		Hash:         mockedHash{nodeInfos: nodeInfos},
		ClosestNodes: closestNodes,
		Found:        true,
	})

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.nodeID,
		Values: []dht.NodeAddr{
			nodeInfos[0].Addr,
			nodeInfos[1].Addr,
			nodeInfos[2].Addr,
		},
		Nodes: nodeInfosFromNodes(closestNodes...),
		Token: &expectedToken,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_get_peers__nodes(t *testing.T) {
	t.Parallel()

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
	nodes := dht.CompactIPv4NodeInfo{
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
		dht.RandomNodeInfo(4),
	}
	peers := []ktable.Node{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
		mockedPeer{nodes[2]},
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())
	mocks.table.On("GetHashOrClosestNodes", infoHash).Return(ktable.GetHashOrClosestNodesResult{
		ClosestNodes: peers,
	})

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, &expectedToken, ret.Token)
	ret.Token = nil
	assert.Equal(t, dht.Return{
		ID:     mocks.nodeID,
		Values: nil,
		Nodes:  nodes,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_get_peers__missing_info_hash(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
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

func TestResponder_get_peers__ipv6_values(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
	infoHash := protocol.RandomNodeID()
	msg := dht.RecvMsg{
		From: mocks.sender6.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "get_peers",
			A: &dht.MsgArgs{
				ID:       mocks.sender6.ID,
				InfoHash: infoHash,
			},
		},
	}
	nodeInfos := []dht.NodeInfo{
		dht.RandomNodeInfo(16),
		dht.RandomNodeInfo(16),
	}
	closestNodes := []ktable.Node{
		mockedPeer{dht.RandomNodeInfo(16)},
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender6.ID, mocks.sender6.Addr.ToAddrPort().Addr())
	mocks.table6.On("GetHashOrClosestNodes", infoHash).Return(ktable.GetHashOrClosestNodesResult{
		Hash:         mockedHash{nodeInfos: nodeInfos},
		ClosestNodes: closestNodes,
		Found:        true,
	})

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, &expectedToken, ret.Token)
	ret.Token = nil
	assert.Equal(t, dht.Return{
		ID: mocks.nodeID,
		Values: []dht.NodeAddr{
			nodeInfos[0].Addr,
			nodeInfos[1].Addr,
		},
		Nodes6: nodeInfosFromNodes(closestNodes...),
	}, ret)
	assert.NoError(t, err)
	mocks.table.AssertNotCalled(t, "GetHashOrClosestNodes")
}

func TestResponder_get_peers__ipv4_want_n6(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
	infoHash := protocol.RandomNodeID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "get_peers",
			A: &dht.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
				Want:     []dht.Want{dht.WantNodes6},
			},
		},
	}
	nodes6 := dht.CompactIPv6NodeInfo{
		dht.RandomNodeInfo(16),
	}
	peers6 := []ktable.Node{
		mockedPeer{nodes6[0]},
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender.ID, mocks.sender.Addr.ToAddrPort().Addr())

	// This will be called because the request is from an IPv4 address, but Nodes should not be populated
	mocks.table.On("GetHashOrClosestNodes", infoHash).Return(ktable.GetHashOrClosestNodesResult{})

	mocks.table6.On("GetHashOrClosestNodes", infoHash).Return(ktable.GetHashOrClosestNodesResult{
		ClosestNodes: peers6,
	})

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, &expectedToken, ret.Token)
	ret.Token = nil
	assert.Equal(t, dht.Return{
		ID:     mocks.nodeID,
		Values: nil,
		Nodes:  nil,
		Nodes6: nodes6,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_announce_peer__implied_port(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
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
			Addr: mocks.sender.Addr.ToAddrPort(),
		}},
	}).Return(nil)

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.nodeID,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_announce_peer__specified_port(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
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
			Addr: netip.AddrPortFrom(mocks.sender.Addr.ToAddrPort().Addr(), uint16(port)),
		}},
	}).Return(nil)

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.nodeID,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_announce_peer__ipv6(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
	infoHash := protocol.RandomNodeID()
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender6.ID, mocks.sender6.Addr.ToAddrPort().Addr())
	msg := dht.RecvMsg{
		From: mocks.sender6.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "announce_peer",
			A: &dht.MsgArgs{
				ID:          mocks.sender6.ID,
				InfoHash:    infoHash,
				Token:       expectedToken,
				ImpliedPort: true,
			},
		},
	}

	mocks.table6.On("BatchCommand", ktable.PutHash{
		ID: infoHash,
		Peers: []ktable.HashPeer{{
			Addr: mocks.sender6.Addr.ToAddrPort(),
		}},
	}).Return(nil)

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, dht.Return{
		ID: mocks.nodeID,
	}, ret)
	assert.NoError(t, err)
	mocks.table.AssertNotCalled(t, "BatchCommand")
}

func TestResponder_sample_infohashes(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
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
	peers := []ktable.Node{
		mockedPeer{nodes[0]},
		mockedPeer{nodes[1]},
		mockedPeer{nodes[2]},
	}
	num := int64(123)
	mocks.table.On("SampleHashesAndNodes").Return(ktable.SampleHashesAndNodesResult{
		Hashes:      infoHashes,
		Nodes:       peers,
		TotalHashes: int(num),
	})

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, &mocks.responder.sampleInfoHashesInterval, ret.Interval)
	assert.Equal(t, dht.Return{
		ID:    mocks.nodeID,
		Nodes: nodes,
		Bep51Return: dht.Bep51Return{
			Samples: &dht.CompactInfohashes{
				infoHashes[0].ID(),
				infoHashes[1].ID(),
				infoHashes[2].ID(),
			},
			Num:      &num,
			Interval: ret.Interval,
		},
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_sample_infohashes__ipv6(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender6.Addr.ToAddrPort(),
		Msg: dht.Msg{
			Q: "sample_infohashes",
			A: &dht.MsgArgs{
				ID: mocks.sender6.ID,
			},
		},
	}
	infoHashes := []ktable.Hash{
		mockedHash{id: protocol.RandomNodeID()},
	}
	nodes := dht.CompactIPv6NodeInfo{
		dht.RandomNodeInfo(16),
	}
	peers := []ktable.Node{
		mockedPeer{nodes[0]},
	}
	num := int64(123)
	mocks.table6.On("SampleHashesAndNodes").Return(ktable.SampleHashesAndNodesResult{
		Hashes:      infoHashes,
		Nodes:       peers,
		TotalHashes: int(num),
	})

	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, &mocks.responder.sampleInfoHashesInterval, ret.Interval)
	assert.Equal(t, dht.Return{
		ID:     mocks.nodeID,
		Nodes6: nodes,
		Bep51Return: dht.Bep51Return{
			Samples:  &dht.CompactInfohashes{infoHashes[0].ID()},
			Num:      &num,
			Interval: ret.Interval,
		},
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_unknown_method(t *testing.T) {
	t.Parallel()

	mocks := newTestResponderMocks(t)
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
