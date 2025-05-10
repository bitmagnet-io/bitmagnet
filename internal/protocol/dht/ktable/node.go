package ktable

import (
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type nodeKeyspace struct {
	keyspace[netip.AddrPort, NodeOption, Node, *node]
}

func (k *nodeKeyspace) getLastRespondedBefore(t time.Time) []Node {
	var peers []Node

	for _, it := range k.items {
		if it.lastRespondedAt.Before(t) {
			peers = append(peers, it)
		}
	}

	return peers
}

func (k *nodeKeyspace) getCandidatesForSampleInfoHashes(n int) []*node {
	//nolint:prealloc
	var candidates []*node

	for _, it := range k.items {
		if !it.IsSampleInfoHashesCandidate() {
			continue
		}

		candidates = append(candidates, it)
		if len(candidates) == n {
			break
		}
	}

	return candidates
}

type Node interface {
	keyspaceItem
	Addr() netip.AddrPort
	Time() time.Time
	Dropped() bool
	IsSampleInfoHashesCandidate() bool
}

type nodeBase struct {
	id   ID
	addr netip.AddrPort
}

func NewNode(id ID, addr netip.AddrPort) Node {
	return nodeBase{
		id:   id,
		addr: addr,
	}
}

func (p nodeBase) ID() protocol.ID {
	return p.id
}

func (p nodeBase) Addr() netip.AddrPort {
	return p.addr
}

func (nodeBase) Time() time.Time {
	return time.Time{}
}

func (nodeBase) Dropped() bool {
	return false
}

func (nodeBase) IsSampleInfoHashesCandidate() bool {
	return true
}

type node struct {
	nodeBase
	discoveredAt    time.Time
	lastRespondedAt time.Time
	dropReason      error

	bep51Support             protocolSupport
	sampledNum               int
	lastDiscoveredNum        int
	totalNum                 int
	nextSampleInfoHashesTime time.Time

	reverseMap *reverseMap
}

type NodeOption interface {
	apply(*node)
}

type nodeOption struct {
	fn func(*node)
}

func (p nodeOption) apply(peer *node) {
	p.fn(peer)
}

type protocolSupport int

const (
	protocolSupportUnknown protocolSupport = iota
	protocolSupportYes
	protocolSupportNo
)

var _ keyspaceItemPrivate[netip.AddrPort, NodeOption, Node] = (*node)(nil)

func (n *node) update(addr netip.AddrPort) {
	if n.addr != addr {
		n.reverseMap.dropAddr(n.addr.Addr())
	}

	n.addr = addr
	n.reverseMap.putAddrPeerID(addr.Addr(), n.id)
}

func (n *node) apply(option NodeOption) {
	option.apply(n)
}

func (n *node) drop(reason error) {
	n.dropReason = reason
	n.reverseMap.dropAddr(n.addr.Addr())
}

func (n *node) public() Node {
	return n
}

func (n *node) Time() time.Time {
	return n.lastRespondedAt
}

func (n *node) Dropped() bool {
	return n.dropReason != nil
}

func (n *node) IsSampleInfoHashesCandidate() bool {
	now := time.Now()
	threshold := now.Add(-(5 * time.Second))

	return n.bep51Support != protocolSupportNo &&
		n.nextSampleInfoHashesTime.Before(now) &&
		n.lastRespondedAt.Before(threshold)
}

func (n *node) olderThan(t time.Time) bool {
	return n.lastRespondedAt.Before(t)
}

func NodeResponded() NodeOption {
	return nodeOption{
		fn: func(n *node) {
			n.lastRespondedAt = time.Now()
		},
	}
}

func NodeBep51Support(supported bool) NodeOption {
	return nodeOption{
		fn: func(n *node) {
			s := protocolSupportNo
			if supported {
				s = protocolSupportYes
			}
			n.bep51Support = s
		},
	}
}

func NodeSampleInfoHashesRes(discoveredNum int, totalNum int, nextSampleTime time.Time) NodeOption {
	return nodeOption{
		fn: func(n *node) {
			n.sampledNum += discoveredNum
			n.lastDiscoveredNum = discoveredNum
			n.totalNum += totalNum
			// a crude way of deprioritizing nodes that gave us no new samples:
			if discoveredNum == 0 {
				now := time.Now()
				if nextSampleTime.Before(now) {
					nextSampleTime = now
				}
				nextSampleTime = nextSampleTime.Add(5 * time.Minute)
			}
			n.nextSampleInfoHashesTime = nextSampleTime
		},
	}
}
