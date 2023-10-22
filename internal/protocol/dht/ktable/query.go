package ktable

import (
	"net/netip"
	"sort"
	"time"
)

type Query[T any] interface {
	execReturn(*table) T
}

var _ Query[[]Peer] = GetClosestPeers{}

type GetClosestPeers struct {
	ID ID
}

func (c GetClosestPeers) execReturn(t *table) []Peer {
	return t.peers.getClosest(c.ID)
}

var _ Query[[]Peer] = GetOldestPeers{}

type GetOldestPeers struct {
	Cutoff time.Time
	N      int
}

func (c GetOldestPeers) execReturn(t *table) []Peer {
	peers := t.peers.getLastRespondedBefore(c.Cutoff)
	sort.Slice(peers, func(i, j int) bool {
		return peers[i].Time().Before(peers[j].Time())
	})
	if c.N > 0 && len(peers) > c.N {
		peers = peers[:c.N]
	}
	return peers
}

var _ Query[[]netip.Addr] = FilterKnownAddrs{}

type FilterKnownAddrs struct {
	Addrs []netip.Addr
}

func (c FilterKnownAddrs) execReturn(t *table) []netip.Addr {
	var unknown []netip.Addr
	for _, addr := range c.Addrs {
		if _, ok := t.addrs.addrs[addr.String()]; !ok {
			unknown = append(unknown, addr)
		}
	}
	return unknown
}

var _ Query[[]Peer] = GetPeersForSampleInfoHashes{}

type GetPeersForSampleInfoHashes struct {
	N int
}

func (c GetPeersForSampleInfoHashes) execReturn(t *table) []Peer {
	peers := make([]Peer, 0, c.N)
	for _, p := range t.peers.getCandidatesForSampleInfoHashes(c.N) {
		peers = append(peers, p)
		if len(peers) >= c.N {
			break
		}
	}
	return peers
}

var _ Query[ID] = GeneratePeerID{}

type GeneratePeerID struct{}

func (c GeneratePeerID) execReturn(t *table) ID {
	return t.peers.generateRandomID()
}

var _ Query[GetHashOrClosestPeersResult] = GetHashOrClosestPeers{}

type GetHashOrClosestPeers struct {
	ID ID
}

func (c GetHashOrClosestPeers) execReturn(t *table) GetHashOrClosestPeersResult {
	h, ok := t.hashes.items[c.ID]
	if ok {
		return GetHashOrClosestPeersResult{
			Hash:  h,
			Found: true,
		}
	}
	closestPeers := t.peers.getClosest(c.ID)
	return GetHashOrClosestPeersResult{
		ClosestPeers: closestPeers,
	}
}

var _ Query[SampleHashesAndPeersResult] = SampleHashesAndPeers{}

type SampleHashesAndPeers struct{}

func (c SampleHashesAndPeers) execReturn(t *table) SampleHashesAndPeersResult {
	nHashes := 8
	hashes := t.hashes.getRandom(nHashes)
	nPeers := 8 + (nHashes - len(hashes))
	peers := t.peers.getRandom(nPeers)
	totalHashes := t.hashes.count()
	return SampleHashesAndPeersResult{
		Hashes:      hashes,
		Peers:       peers,
		TotalHashes: totalHashes,
	}
}
