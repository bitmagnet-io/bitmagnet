package ktable

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"
	"go.uber.org/fx"
	"net/netip"
	"sync"
	"time"
)

type ID = protocol.ID

type Params struct {
	fx.In
	NodeID ID `name:"dht_node_id"`
}

type Result struct {
	fx.Out
	Table Table
}

func New(p Params) Result {
	return Result{
		Table: newTable(p.NodeID, 80, 80),
	}
}

type TableOrigin interface {
	Origin() ID
}

type TableCommand interface {
	PutNode(id ID, addr netip.AddrPort, options ...NodeOption) btree.PutResult
	DropNode(id ID, reason error) bool
	PutHash(id ID, peers []HashPeer, options ...HashOption) btree.PutResult
}

type TableQuery interface {
	GetClosestNodes(id ID) []Node
	GetOldestNodes(cutoff time.Time, n int) []Node
	GetNodesForSampleInfoHashes(n int) []Node
	FilterKnownAddrs(addrs []netip.Addr) []netip.Addr
  // GenerateNodeID generates a pseudo-random node ID, biased towards the emptiest node buckets.
	GenerateNodeID() ID
	GetHashOrClosestNodes(id ID) GetHashOrClosestNodesResult
  // SampleHashesAndNodes returns a random sample of up to 8 hashes and nodes, and the total hashes count.
	SampleHashesAndNodes() SampleHashesAndNodesResult
}

type Table interface {
	TableOrigin
	TableCommand
	TableQuery
	BatchCommand(commands ...Command)
	Stats() Stats
}

type GetHashOrClosestNodesResult struct {
	Hash         Hash
	ClosestNodes []Node
	Found        bool
}

type SampleHashesAndNodesResult struct {
	Hashes      []Hash
	Nodes       []Node
	TotalHashes int
}

func newTable(origin ID, peersK int, hashesK int) *table {
	rm := &reverseMap{addrs: make(map[string]*infoForAddr)}
	return &table{
		origin:  origin,
		peersK:  peersK,
		hashesK: hashesK,
		nodes: nodeBucketRoot{
			bucketRoot: newBucketRoot[netip.AddrPort, NodeOption, Node, *node](
				origin,
				peersK,
				func(id ID, addr netip.AddrPort) *node {
					return &node{
						nodeBase: nodeBase{
							id:   id,
							addr: addr,
						},
						discoveredAt: time.Now(),
						reverseMap:   rm,
					}
				},
			),
		},
		hashes: hashBucketRoot{
			bucketRoot: newBucketRoot[[]HashPeer, HashOption, Hash, *hash](
				origin,
				hashesK,
				func(id ID, peers []HashPeer) *hash {
					peersMap := make(map[string]HashPeer, len(peers))
					for _, p := range peers {
						peersMap[p.Addr.Addr().String()] = p
						rm.putAddrHashes(p.Addr.Addr(), id)
					}
					return &hash{
						id:           id,
						peers:        peersMap,
						discoveredAt: time.Now(),
						reverseMap:   rm,
					}
				},
			),
		},
		addrs: rm,
	}
}

type table struct {
	mutex   sync.RWMutex
	origin  ID
	peersK  int
	hashesK int
	nodes   nodeBucketRoot
	hashes  hashBucketRoot
	addrs   *reverseMap
}

func (t *table) Origin() ID {
	return t.origin
}

func (t *table) BatchCommand(commands ...Command) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, command := range commands {
		command.exec(t)
	}
}

func (t *table) PutNode(id ID, addr netip.AddrPort, options ...NodeOption) btree.PutResult {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return PutNode{
		ID:      id,
		Addr:    addr,
		Options: options,
	}.execReturn(t)
}

func (t *table) DropNode(id ID, reason error) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return DropNode{
		ID:     id,
		Reason: reason,
	}.execReturn(t)
}

func (t *table) PutHash(id ID, peers []HashPeer, options ...HashOption) btree.PutResult {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return PutHash{
		ID:      id,
		Peers:   peers,
		Options: options,
	}.execReturn(t)
}

func (t *table) GetClosestNodes(id ID) []Node {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return GetClosestPeers{
		ID: id,
	}.execReturn(t)
}

func (t *table) GetOldestNodes(cutoff time.Time, n int) []Node {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return GetOldestPeers{
		Cutoff: cutoff,
		N:      n,
	}.execReturn(t)
}

func (t *table) GetNodesForSampleInfoHashes(n int) []Node {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return GetNodesForSampleInfoHashes{
		N: n,
	}.execReturn(t)
}

func (t *table) FilterKnownAddrs(addrs []netip.Addr) []netip.Addr {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return FilterKnownAddrs{
		Addrs: addrs,
	}.execReturn(t)
}

func (t *table) GenerateNodeID() ID {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return GenerateNodeID{}.execReturn(t)
}

func (t *table) GetHashOrClosestNodes(id ID) GetHashOrClosestNodesResult {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return GetHashOrClosestNodes{
		ID: id,
	}.execReturn(t)
}

func (t *table) SampleHashesAndNodes() SampleHashesAndNodesResult {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return SampleHashesAndNodes{}.execReturn(t)
}
