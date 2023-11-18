package ktable

import (
	"net/netip"
	"sort"
	"time"
)

type Query[T any] interface {
	execReturn(*table) T
}

var _ Query[[]Node] = GetClosestPeers{}

type GetClosestPeers struct {
	ID ID
}

func (c GetClosestPeers) execReturn(t *table) []Node {
	return t.nodes.getClosest(c.ID)
}

var _ Query[[]Node] = GetOldestPeers{}

type GetOldestPeers struct {
	Cutoff time.Time
	N      int
}

func (c GetOldestPeers) execReturn(t *table) []Node {
	peers := t.nodes.getLastRespondedBefore(c.Cutoff)
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

var _ Query[[]Node] = GetNodesForSampleInfoHashes{}

type GetNodesForSampleInfoHashes struct {
	N int
}

func (c GetNodesForSampleInfoHashes) execReturn(t *table) []Node {
	peers := make([]Node, 0, c.N)
	for _, p := range t.nodes.getCandidatesForSampleInfoHashes(c.N) {
		peers = append(peers, p)
		if len(peers) >= c.N {
			break
		}
	}
	return peers
}

var _ Query[GetHashOrClosestNodesResult] = GetHashOrClosestNodes{}

type GetHashOrClosestNodes struct {
	ID ID
}

func (c GetHashOrClosestNodes) execReturn(t *table) GetHashOrClosestNodesResult {
	h, ok := t.hashes.get(c.ID)
	if ok {
		return GetHashOrClosestNodesResult{
			Hash:  h,
			Found: true,
		}
	}
	closestNodes := t.nodes.getClosest(c.ID)
	return GetHashOrClosestNodesResult{
		ClosestNodes: closestNodes,
	}
}

var _ Query[SampleHashesAndNodesResult] = SampleHashesAndNodes{}

type SampleHashesAndNodes struct{}

func (c SampleHashesAndNodes) execReturn(t *table) SampleHashesAndNodesResult {
	nHashes := 20
	hashes := t.hashes.getRandom(nHashes)
	nNodes := 20 + (nHashes - len(hashes))
	nodes := t.nodes.getRandom(nNodes)
	totalHashes := t.hashes.count()
	return SampleHashesAndNodesResult{
		Hashes:      hashes,
		Nodes:       nodes,
		TotalHashes: totalHashes,
	}
}
