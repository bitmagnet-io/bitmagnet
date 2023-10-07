package routingtable

import (
	"context"
	"errors"
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/lru"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"sort"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Config dht.Config
}

type Result struct {
	fx.Out
	Table Table
}

func New(p Params) Result {
	findNode, _ := lru.New[krpc.ID, krpc.NodeAddr](10000)
	getPeers, _ := lru.New[krpc.ID, peersForHash](10000)
	return Result{
		Table: &table{
			mutex:                 &sync.RWMutex{},
			semaphore:             semaphore.NewWeighted(int64(p.Config.Routing.MaxConcurrency)),
			crawlerNodes:          make(map[string]*crawlerNode, p.Config.Routing.MaxPeers),
			maxNodes:              p.Config.Routing.MaxPeers,
			maxConcurrencyPerPeer: p.Config.Routing.MaxConcurrencyPerPeer,
			findNode:              findNode,
			getPeers:              getPeers,
		},
	}
}

type Table interface {
	ReceiveNodeInfo(...krpc.NodeInfo)
	ReceiveNodeAddr(...krpc.NodeAddr)
	ReceivePeersForHash(hash krpc.ID, addrs ...krpc.NodeAddr)
	FindNode(krpc.ID) krpc.CompactIPv4NodeInfo
	GetPeers(hash krpc.ID) ([]krpc.NodeAddr, krpc.CompactIPv4NodeInfo)
	SampleInfoHashes() (krpc.CompactInfohashes, krpc.CompactIPv4NodeInfo, int64)
	WithPeer(context.Context, krpc.NodeAddr, func(ctx context.Context) error) error
	// TryEachNode tries to execute the given function on each node in the routing table, until the function returns nil.
	TryEachNode(context.Context, func(ctx context.Context, peer PeerInfo) error) error
}

type table struct {
	mutex                 *sync.RWMutex
	semaphore             *semaphore.Weighted
	crawlerNodes          map[string]*crawlerNode
	maxNodes              uint
	maxConcurrencyPerPeer uint
	// bep5, bep51 stores:
	// findNode a map of node IDs to node addresses
	findNode *lru.Cache[krpc.ID, krpc.NodeAddr]
	// getPeers a map of infohashes to a set of node addresses
	getPeers *lru.Cache[krpc.ID, peersForHash]
}

type peersForHash struct {
	mutex *sync.RWMutex
	addrs *map[string]krpc.NodeAddr
}

func newPeersForHash() peersForHash {
	a := make(map[string]krpc.NodeAddr)
	return peersForHash{
		mutex: &sync.RWMutex{},
		addrs: &a,
	}
}

func (p peersForHash) receiveAddrs(addrs ...krpc.NodeAddr) {
	p.mutex.Lock()
	for _, addr := range addrs {
		(*p.addrs)[addr.String()] = addr
	}
	p.mutex.Unlock()
}

func (p peersForHash) getAddrs() []krpc.NodeAddr {
	p.mutex.RLock()
	addrs := make([]krpc.NodeAddr, 0, 8)
	for _, addr := range *p.addrs {
		addrs = append(addrs, addr)
		if len(addrs) == 8 {
			break
		}
	}
	p.mutex.RUnlock()
	return addrs
}

func (t *table) ReceiveNodeInfo(nodes ...krpc.NodeInfo) {
	addrs := make([]krpc.NodeAddr, 0, len(nodes))
	for _, node := range nodes {
		addrs = append(addrs, node.Addr)
		t.findNode.Add(node.ID, node.Addr)
	}
	t.ReceiveNodeAddr(addrs...)
}

func (t *table) ReceiveNodeAddr(peers ...krpc.NodeAddr) {
	if len(peers) == 0 {
		return
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	now := time.Now()
	peersToAdd := make([]krpc.NodeAddr, 0, len(peers))
	for _, nodeAddr := range peers {
		if nodeAddr.Port == 0 {
			continue
		}
		if existingPeer, ok := t.crawlerNodes[nodeAddr.String()]; !ok {
			peersToAdd = append(peersToAdd, nodeAddr)
		} else {
			existingPeer.lastDiscoveredAt = now
			existingPeer.discoveredCount++
		}
	}
	nPeersToAdd := len(peersToAdd)
	maxPeersToAdd := int(t.maxNodes) - len(t.crawlerNodes)
	if nPeersToAdd > maxPeersToAdd {
		nEvicted := t.tryEvictPeersLocked(nPeersToAdd - maxPeersToAdd)
		nPeersToAdd = maxPeersToAdd + nEvicted
	}
	for i := 0; i < nPeersToAdd; i++ {
		nodeAddr := peersToAdd[i]
		newPeer := t.newPeer(nodeAddr)
		t.crawlerNodes[nodeAddr.String()] = newPeer
	}
}

func (t *table) newPeer(nodeAddr krpc.NodeAddr) *crawlerNode {
	now := time.Now()
	return &crawlerNode{
		peerSemaphore:    semaphore.NewWeighted(int64(t.maxConcurrencyPerPeer)),
		tableSemaphore:   t.semaphore,
		nodeAddr:         nodeAddr,
		discoveredAt:     now,
		lastDiscoveredAt: now,
		discoveredCount:  1,
	}
}

// tryEvictPeersLocked tries to evict n crawlerNodes from the routing table based on the oldest crawlerNodes to receive a successful response that can be unlocked.
// Will only consider crawlerNodes with at least one response or error.
// Assumes that the caller has already locked the routing table.
// todo: Improve eviction strategy.
func (t *table) tryEvictPeersLocked(n int) int {
	maxCandidates := n * 5
	if maxCandidates > len(t.crawlerNodes) {
		maxCandidates = len(t.crawlerNodes)
	}
	candidates := make([]*crawlerNode, 0, maxCandidates)
	for _, peer := range t.crawlerNodes {
		if peer.responseCount+peer.errorCount > 0 && peer.peerSemaphore.TryAcquire(int64(t.maxConcurrencyPerPeer)) {
			candidates = append(candidates, peer)
			if len(candidates) >= maxCandidates {
				break
			}
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].lastRespondedAt.Before(candidates[j].lastRespondedAt)
	})
	nEvicted := 0
	for i := 0; i < len(candidates); i++ {
		if nEvicted < n {
			candidates[i].evicted = true
			delete(t.crawlerNodes, candidates[i].nodeAddr.String())
			nEvicted++
		}
		candidates[i].peerSemaphore.Release(int64(t.maxConcurrencyPerPeer))
	}
	return nEvicted
}

func (t *table) ReceivePeersForHash(hash krpc.ID, addrs ...krpc.NodeAddr) {
	current, ok := t.getPeers.Peek(hash)
	if !ok {
		current = newPeersForHash()
	}
	current.receiveAddrs(addrs...)
	t.getPeers.Add(hash, current)
}

func (t *table) FindNode(id krpc.ID) krpc.CompactIPv4NodeInfo {
	addr, ok := t.findNode.Get(id)
	if ok {
		return []krpc.NodeInfo{
			{
				ID:   id,
				Addr: addr,
			},
		}
	}
	sample := t.findNode.Sample(8)
	addrs := make([]krpc.NodeInfo, 0, len(sample))
	for _, e := range sample {
		addrs = append(addrs, krpc.NodeInfo{
			ID:   e.Key,
			Addr: e.Value,
		})
	}
	return addrs
}

func (t *table) GetPeers(hash krpc.ID) ([]krpc.NodeAddr, krpc.CompactIPv4NodeInfo) {
	pfh, ok := t.getPeers.Get(hash)
	if ok {
		return pfh.getAddrs(), nil
	}
	nodesSample := t.findNode.Sample(8)
	nodes := make(krpc.CompactIPv4NodeInfo, 0, len(nodesSample))
	for _, e := range nodesSample {
		nodes = append(nodes, krpc.NodeInfo{
			ID:   e.Key,
			Addr: e.Value,
		})
	}
	return nil, nodes
}

func (t *table) AnnouncePeer(hash krpc.ID, node krpc.NodeInfo) {
	t.ReceivePeersForHash(hash, node.Addr)
}

func (t *table) SampleInfoHashes() (krpc.CompactInfohashes, krpc.CompactIPv4NodeInfo, int64) {
	hashesSample := t.getPeers.Sample(8)
	hashes := make(krpc.CompactInfohashes, 0, len(hashesSample))
	for _, e := range hashesSample {
		hashes = append(hashes, e.Key)
	}
	nodesSample := t.findNode.Sample(8)
	nodes := make(krpc.CompactIPv4NodeInfo, 0, len(nodesSample))
	for _, e := range nodesSample {
		nodes = append(nodes, krpc.NodeInfo{
			ID:   e.Key,
			Addr: e.Value,
		})
	}
	return hashes, nodes, int64(t.getPeers.Len())
}

func (t *table) TryEachNode(
	ctx context.Context,
	fn func(ctx context.Context, peer PeerInfo) error,
) error {
	t.mutex.RLock()
	peersToTry := make([]*crawlerNode, 0, len(t.crawlerNodes))
	for _, peer := range t.crawlerNodes {
		peersToTry = append(peersToTry, peer)
	}
	t.mutex.RUnlock()
	var err error
	for _, peer := range peersToTry {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		ok, thisErr := t.tryPeerLocked(ctx, peer, fn)
		if ok {
			if thisErr != nil {
				err = thisErr
			} else {
				return nil
			}
		}
	}
	if err == nil {
		err = fmt.Errorf("no crawlerNode available (total %d)", len(peersToTry))
	}
	return err
}

func (t *table) tryPeerLocked(
	ctx context.Context,
	peer *crawlerNode,
	fn func(ctx context.Context, peer PeerInfo) error,
) (bool, error) {
	ok := false
	errChan := make(chan error)
	go (func() {
		_, err := peer.tryWithLock(ctx, func(ctx context.Context) error {
			ok = true
			return fn(ctx, peer)
		})
		errChan <- err
		close(errChan)
	})()
	select {
	case <-ctx.Done():
		return ok, ctx.Err()
	case err := <-errChan:
		return ok, err
	}
}

func (t *table) WithPeer(ctx context.Context, info krpc.NodeAddr, fn func(ctx context.Context) error) error {
	t.mutex.Lock()
	thePeer, ok := t.crawlerNodes[info.String()]
	if !ok {
		diff := len(t.crawlerNodes) - int(t.maxNodes)
		if diff > 0 {
			_ = t.tryEvictPeersLocked(diff)
		}
		thePeer = t.newPeer(info)
		t.crawlerNodes[info.String()] = thePeer
	}
	t.mutex.Unlock()
	return thePeer.WithLock(ctx, fn)
}

type PeerInfo interface {
	Addr() krpc.NodeAddr
	ResponseCount() uint
	LastRespondedAt() time.Time
}

type Peer interface {
	PeerInfo
	WithLock(ctx context.Context, fn func(ctx context.Context) error) error
}

type crawlerNode struct {
	peerSemaphore    *semaphore.Weighted
	tableSemaphore   *semaphore.Weighted
	nodeAddr         krpc.NodeAddr
	discoveredAt     time.Time
	discoveredCount  uint
	lastDiscoveredAt time.Time
	responseCount    uint
	lastRespondedAt  time.Time
	errorCount       uint
	lastErroredAt    time.Time
	lastError        error
	evicted          bool
}

func (p *crawlerNode) Addr() krpc.NodeAddr {
	return p.nodeAddr
}

func (p *crawlerNode) ResponseCount() uint {
	return p.responseCount
}

func (p *crawlerNode) LastRespondedAt() time.Time {
	return p.lastRespondedAt
}

func (p *crawlerNode) WithLock(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := p.peerSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer p.peerSemaphore.Release(1)
	if err := p.tableSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer p.tableSemaphore.Release(1)
	return p.doLocked(ctx, fn)
}

func (p *crawlerNode) tryWithLock(ctx context.Context, fn func(ctx context.Context) error) (bool, error) {
	if ok := p.peerSemaphore.TryAcquire(1); !ok {
		return false, nil
	}
	defer p.peerSemaphore.Release(1)
	if ok := p.tableSemaphore.TryAcquire(1); !ok {
		return false, nil
	}
	defer p.tableSemaphore.Release(1)
	err := p.doLocked(ctx, fn)
	return true, err
}

func (p *crawlerNode) doLocked(ctx context.Context, fn func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if p.evicted {
		return errors.New("crawlerNode evicted")
	}
	err := fn(ctx)
	if err != nil {
		p.errorCount++
		p.lastErroredAt = time.Now()
		p.lastError = err
	} else {
		p.responseCount++
		p.lastRespondedAt = time.Now()
	}
	return err
}
