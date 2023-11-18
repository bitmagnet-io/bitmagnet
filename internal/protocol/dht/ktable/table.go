package ktable

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"
	"net/netip"
	"sync"
	"time"
)

type ID = protocol.ID

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
	GetHashOrClosestNodes(id ID) GetHashOrClosestNodesResult
	// SampleHashesAndNodes returns a random sample of up to 8 hashes and nodes, and the total hashes count.
	SampleHashesAndNodes() SampleHashesAndNodesResult
}

type Table interface {
	TableOrigin
	TableCommand
	TableQuery
	BatchCommand(commands ...Command)
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

type table struct {
	mutex   sync.RWMutex
	origin  ID
	nodesK  int
	hashesK int
	nodes   nodeKeyspace
	hashes  hashKeyspace
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
