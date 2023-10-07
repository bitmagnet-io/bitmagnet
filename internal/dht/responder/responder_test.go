package responder

import (
	"context"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/mocks/dht/routingtable"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testResponderMocks struct {
	peerID    krpc.ID
	table     *routingtable.Table
	responder responder
	sender    krpc.NodeInfo
}

func newTestResponderMocks(t *testing.T) testResponderMocks {
	peerID := dht.RandomPeerID()
	tableMock := routingtable.NewTable(t)
	return testResponderMocks{
		peerID: peerID,
		table:  tableMock,
		responder: responder{
			peerID:                   peerID,
			table:                    tableMock,
			sampleInfoHashesInterval: 20,
		},
		sender: krpc.RandomNodeInfo(4),
	}
}

func TestResponder_ping(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "ping",
			A: &krpc.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{ID: mocks.peerID}, ret)
	assert.NoError(t, err)
}

func TestResponder_ping__missing_args(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "ping",
		},
	}
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_find_node(t *testing.T) {
	mocks := newTestResponderMocks(t)
	target := dht.RandomPeerID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "find_node",
			A: &krpc.MsgArgs{
				ID:     mocks.sender.ID,
				Target: target,
			},
		},
	}
	nodes := krpc.CompactIPv4NodeInfo{
		krpc.RandomNodeInfo(4),
		krpc.RandomNodeInfo(4),
		krpc.RandomNodeInfo(4),
	}
	mocks.table.On("FindNode", target).Return(nodes)
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{
		ID:    mocks.peerID,
		Nodes: nodes,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_find_node__missing_target(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "find_node",
			A: &krpc.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, ErrMissingArguments, err)
}

func TestResponder_get_peers__values(t *testing.T) {
	mocks := newTestResponderMocks(t)
	infoHash := dht.RandomPeerID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "get_peers",
			A: &krpc.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
			},
		},
	}
	values := []krpc.NodeAddr{
		krpc.RandomNodeInfo(4).Addr,
		krpc.RandomNodeInfo(4).Addr,
		krpc.RandomNodeInfo(4).Addr,
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender)
	mocks.table.On("GetPeers", infoHash).Return(values, nil)
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{
		ID:     mocks.peerID,
		Values: values,
		Nodes:  nil,
		Token:  &expectedToken,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_get_peers__nodes(t *testing.T) {
	mocks := newTestResponderMocks(t)
	infoHash := dht.RandomPeerID()
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "get_peers",
			A: &krpc.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
			},
		},
	}
	nodes := krpc.CompactIPv4NodeInfo{
		krpc.RandomNodeInfo(4),
		krpc.RandomNodeInfo(4),
		krpc.RandomNodeInfo(4),
	}
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender)
	mocks.table.On("GetPeers", infoHash).Return(nil, nodes)
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{
		ID:     mocks.peerID,
		Values: nil,
		Nodes:  nodes,
		Token:  &expectedToken,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_announce_peer__implied_port(t *testing.T) {
	mocks := newTestResponderMocks(t)
	infoHash := dht.RandomPeerID()
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "announce_peer",
			A: &krpc.MsgArgs{
				ID:          mocks.sender.ID,
				InfoHash:    infoHash,
				Token:       expectedToken,
				ImpliedPort: true,
			},
		},
	}
	mocks.table.On("ReceivePeersForHash", infoHash, mocks.sender.Addr).Return()
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{
		ID: mocks.peerID,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_announce_peer__specified_port(t *testing.T) {
	mocks := newTestResponderMocks(t)
	infoHash := dht.RandomPeerID()
	expectedToken := mocks.responder.announceToken(infoHash, mocks.sender)
	port := krpc.RandomNodeInfo(4).Addr.Port
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "announce_peer",
			A: &krpc.MsgArgs{
				ID:       mocks.sender.ID,
				InfoHash: infoHash,
				Token:    expectedToken,
				Port:     &port,
			},
		},
	}
	mocks.table.On("ReceivePeersForHash", infoHash, krpc.NodeAddr{
		IP:   mocks.sender.Addr.IP,
		Port: port,
	}).Return()
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{
		ID: mocks.peerID,
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_sample_infohashes(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "sample_infohashes",
			A: &krpc.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	infoHashes := krpc.CompactInfohashes{
		dht.RandomPeerID(),
		dht.RandomPeerID(),
		dht.RandomPeerID(),
	}
	nodes := krpc.CompactIPv4NodeInfo{
		krpc.RandomNodeInfo(4),
		krpc.RandomNodeInfo(4),
		krpc.RandomNodeInfo(4),
	}
	num := int64(123)
	mocks.table.On("SampleInfoHashes").Return(infoHashes, nodes, num)
	mocks.table.On("ReceiveNodeInfo", mocks.sender).Return()
	ret, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, krpc.Return{
		ID:    mocks.peerID,
		Nodes: nodes,
		Bep51Return: krpc.Bep51Return{
			Samples:  &infoHashes,
			Num:      &num,
			Interval: &mocks.responder.sampleInfoHashesInterval,
		},
	}, ret)
	assert.NoError(t, err)
}

func TestResponder_unknown_method(t *testing.T) {
	mocks := newTestResponderMocks(t)
	msg := dht.RecvMsg{
		From: mocks.sender.Addr,
		Msg: krpc.Msg{
			Q: "foo",
			A: &krpc.MsgArgs{
				ID: mocks.sender.ID,
			},
		},
	}
	_, err := mocks.responder.Respond(context.Background(), msg)
	assert.Equal(t, err, ErrMethodUnknown)
}
