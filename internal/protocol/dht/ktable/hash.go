package ktable

import (
	"net/netip"
	"time"
)

type hashKeyspace struct {
	keyspace[[]HashPeer, HashOption, Hash, *hash]
}

type Hash interface {
	keyspaceItem
	Peers() []HashPeer
	Dropped() bool
}

type hash struct {
	id           ID
	peers        map[string]HashPeer
	discoveredAt time.Time
	// lastRequestedAt time.Time
	droppedReason error
	reverseMap    *reverseMap
}

type HashPeer struct {
	Addr netip.AddrPort
}

type HashOption interface {
	apply(*hash)
}

var _ keyspaceItemPrivate[[]HashPeer, HashOption, Hash] = (*hash)(nil)

func (h *hash) update(peers []HashPeer) {
	for _, p := range peers {
		h.peers[p.Addr.Addr().String()] = p
		h.reverseMap.putAddrHashes(p.Addr.Addr(), h.id)
	}
}

func (h *hash) apply(option HashOption) {
	option.apply(h)
}

func (h *hash) drop(reason error) {
	h.droppedReason = reason
	for _, addr := range h.peers {
		if info, ok := h.reverseMap.addrs[addr.Addr.Addr().String()]; ok {
			info.dropHashes(h.id)
		}
	}
}

func (h *hash) public() Hash {
	return h
}

func (h *hash) hasPeers() bool {
	return len(h.peers) > 0
}

func (h *hash) ID() ID {
	return h.id
}

func (h *hash) Peers() []HashPeer {
	peers := make([]HashPeer, 0, len(h.peers))
	for _, p := range h.peers {
		peers = append(peers, p)
	}

	return peers
}

func (h *hash) Dropped() bool {
	return h.droppedReason != nil
}
