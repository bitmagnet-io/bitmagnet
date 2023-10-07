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
	"net"
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
	return Result{
		Table: &table{
			mutex:                 &sync.RWMutex{},
			semaphore:             semaphore.NewWeighted(int64(p.Config.Routing.MaxConcurrency)),
			crawlerNodes:          make(map[string]*crawlerNode, p.Config.Routing.MaxPeers),
			maxNodes:              p.Config.Routing.MaxPeers,
			maxConcurrencyPerPeer: p.Config.Routing.MaxConcurrencyPerPeer,
			findNode:              lru.NewExpirable[krpc.ID, krpc.NodeAddr](10000, nil, time.Minute*60),
			getPeers:              lru.NewExpirable[krpc.ID, peersForHash](10000, nil, time.Minute*60),
			goodBadNodes:          lru.NewExpirable[string, bool](10000, nil, time.Minute*15),
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
	WithPeer(context.Context, krpc.NodeAddr, func(context.Context) error) error
	// TryEachNode tries to execute the given function on each node in the routing table, until the function returns nil.
	TryEachNode(context.Context, func(context.Context, krpc.NodeAddr) error) error
}

type table struct {
	mutex                 *sync.RWMutex
	semaphore             *semaphore.Weighted
	crawlerNodes          map[string]*crawlerNode
	maxNodes              uint
	maxConcurrencyPerPeer uint
	// bep5, bep51 stores:
	// findNode an expiring map of node IDs to node addresses
	findNode *lru.Expirable[krpc.ID, krpc.NodeAddr]
	// getPeers an expiring map of infohashes to a set of node addresses
	getPeers *lru.Expirable[krpc.ID, peersForHash]
	// goodBadNodes an expiring map of node IPs to a bool indicating if a success or error was received
	goodBadNodes *lru.Expirable[string, bool]
}

type peersForHash struct {
	table *table
	mutex *sync.RWMutex
	addrs *map[string]krpc.NodeAddr
}

func (t *table) newPeersForHash() peersForHash {
	a := make(map[string]krpc.NodeAddr)
	return peersForHash{
		table: t,
		mutex: &sync.RWMutex{},
		addrs: &a,
	}
}

func (p peersForHash) receiveAddrs(addrs ...krpc.NodeAddr) {
	p.mutex.Lock()
	for _, addr := range addrs {
		if p.table.isBadNode(addr.IP) {
			continue
		}
		(*p.addrs)[addr.String()] = addr
	}
	p.mutex.Unlock()
}

func (p peersForHash) getAddrs() []krpc.NodeAddr {
	target := 8
	p.mutex.RLock()
	sampledFirstPass := make(map[string]struct{}, target)
	addrs := make([]krpc.NodeAddr, 0, target)
	// first try to get known good nodes
	for _, addr := range *p.addrs {
		if p.table.isGoodNode(addr.IP) {
			addrs = append(addrs, addr)
			sampledFirstPass[addr.String()] = struct{}{}
			if len(addrs) == target {
				break
			}
		}
	}
	// if not enough then add nodes that aren't bad
	if len(addrs) < target {
		for _, addr := range *p.addrs {
			if _, ok := sampledFirstPass[addr.String()]; ok {
				continue
			}
			if !p.table.isBadNode(addr.IP) {
				addrs = append(addrs, addr)
				if len(addrs) == target {
					break
				}
			}
		}
	}
	p.mutex.RUnlock()
	return addrs
}

func (t *table) isBadNode(ip net.IP) bool {
	good, ok := t.goodBadNodes.Get(ip.String())
	return ok && !good
}

func (t *table) isGoodNode(ip net.IP) bool {
	good, ok := t.goodBadNodes.Get(ip.String())
	return ok && good
}

func (t *table) ReceiveNodeInfo(nodes ...krpc.NodeInfo) {
	addrs := make([]krpc.NodeAddr, 0, len(nodes))
	for _, node := range nodes {
		if t.isBadNode(node.Addr.IP) {
			continue
		}
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
	peersToAdd := make([]krpc.NodeAddr, 0, len(peers))
	for _, nodeAddr := range peers {
		if nodeAddr.Port == 0 {
			continue
		}
		if t.isBadNode(nodeAddr.IP) {
			continue
		}
		if _, ok := t.crawlerNodes[nodeAddr.String()]; !ok {
			peersToAdd = append(peersToAdd, nodeAddr)
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
	return &crawlerNode{
		table:     t,
		semaphore: semaphore.NewWeighted(int64(t.maxConcurrencyPerPeer)),
		nodeAddr:  nodeAddr,
	}
}

// tryEvictPeersLocked tries to evict n crawlerNodes from the routing table based on known bad nodes.
// Assumes that the caller has already locked the routing table.
func (t *table) tryEvictPeersLocked(n int) int {
	nEvicted := 0
	for _, peer := range t.crawlerNodes {
		if t.isBadNode(peer.nodeAddr.IP) {
			peer.evicted = true
			delete(t.crawlerNodes, peer.nodeAddr.String())
			nEvicted++
			if nEvicted >= n {
				break
			}
		}
	}
	return nEvicted
}

func (t *table) ReceivePeersForHash(hash krpc.ID, addrs ...krpc.NodeAddr) {
	current, ok := t.getPeers.Peek(hash)
	if !ok {
		current = t.newPeersForHash()
	}
	current.receiveAddrs(addrs...)
	t.getPeers.Add(hash, current)
}

func (t *table) sampleBestNodes() krpc.CompactIPv4NodeInfo {
	target := 8
	firstPassSampled := make(map[krpc.ID]struct{}, target)
	// try to return all good nodes
	sample := t.findNode.Sample(8, func(id krpc.ID, addr krpc.NodeAddr) bool {
		firstPassSampled[id] = struct{}{}
		return t.isGoodNode(addr.IP)
	})
	// if not enough then add nodes that aren't bad
	if len(sample) < target {
		sample = append(sample, t.findNode.Sample(target-len(sample), func(id krpc.ID, addr krpc.NodeAddr) bool {
			if _, ok := firstPassSampled[id]; ok {
				return false
			}
			return !t.isBadNode(addr.IP)
		})...)
	}
	addrs := make([]krpc.NodeInfo, 0, len(sample))
	for _, e := range sample {
		addrs = append(addrs, krpc.NodeInfo{
			ID:   e.Key,
			Addr: e.Value,
		})
	}
	return addrs
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
	return t.sampleBestNodes()
}

func (t *table) GetPeers(hash krpc.ID) ([]krpc.NodeAddr, krpc.CompactIPv4NodeInfo) {
	pfh, ok := t.getPeers.Get(hash)
	if ok {
		return pfh.getAddrs(), nil
	}
	return nil, t.sampleBestNodes()
}

func (t *table) AnnouncePeer(hash krpc.ID, node krpc.NodeInfo) {
	t.ReceivePeersForHash(hash, node.Addr)
}

func (t *table) SampleInfoHashes() (krpc.CompactInfohashes, krpc.CompactIPv4NodeInfo, int64) {
	hashesSample := t.getPeers.Sample(8, nil)
	hashes := make(krpc.CompactInfohashes, 0, len(hashesSample))
	for _, e := range hashesSample {
		hashes = append(hashes, e.Key)
	}
	return hashes, t.sampleBestNodes(), int64(t.getPeers.Len())
}

func (t *table) TryEachNode(
	ctx context.Context,
	fn func(context.Context, krpc.NodeAddr) error,
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
	fn func(context.Context, krpc.NodeAddr) error,
) (bool, error) {
	ok := false
	errChan := make(chan error)
	go (func() {
		_, err := peer.tryWithLock(ctx, func(ctx context.Context) error {
			ok = true
			return fn(ctx, peer.nodeAddr)
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

type Peer interface {
	Addr() krpc.NodeAddr
	WithLock(ctx context.Context, fn func(ctx context.Context) error) error
}

type crawlerNode struct {
	table     *table
	semaphore *semaphore.Weighted
	nodeAddr  krpc.NodeAddr
	evicted   bool
}

func (p *crawlerNode) Addr() krpc.NodeAddr {
	return p.nodeAddr
}

func (p *crawlerNode) WithLock(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := p.semaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer p.semaphore.Release(1)
	if err := p.table.semaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer p.table.semaphore.Release(1)
	return p.doLocked(ctx, fn)
}

func (p *crawlerNode) tryWithLock(ctx context.Context, fn func(ctx context.Context) error) (bool, error) {
	if ok := p.semaphore.TryAcquire(1); !ok {
		return false, nil
	}
	defer p.semaphore.Release(1)
	if ok := p.table.semaphore.TryAcquire(1); !ok {
		return false, nil
	}
	defer p.table.semaphore.Release(1)
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
	p.table.goodBadNodes.Add(p.nodeAddr.IP.String(), err == nil)
	return err
}
