package routing

import (
	"context"
	"errors"
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"golang.org/x/sync/semaphore"
	"sort"
	"sync"
	"time"
)

type Config struct {
	MaxPeers              uint
	MaxConcurrency        uint
	MaxConcurrencyPerPeer uint
}

func NewDefaultConfig() Config {
	return Config{
		MaxPeers:              1000,
		MaxConcurrency:        500,
		MaxConcurrencyPerPeer: 5,
	}
}

type Table interface {
	ReceivePeers(peers ...krpc.NodeAddr) (received uint, discarded uint)
	WithPeer(ctx context.Context, info krpc.NodeAddr, f func(ctx context.Context) error) error
	// TryEachPeer tries to execute the given function on each peer in the routing table, until the function returns nil.
	TryEachPeer(ctx context.Context, f func(ctx context.Context, peer PeerInfo) error) error
}

type table struct {
	mutex                 *sync.RWMutex
	semaphore             *semaphore.Weighted
	peers                 map[string]*peer
	maxPeers              uint
	maxConcurrencyPerPeer uint
}

func NewTable(config Config) Table {
	return &table{
		mutex:                 &sync.RWMutex{},
		semaphore:             semaphore.NewWeighted(int64(config.MaxConcurrency)),
		peers:                 make(map[string]*peer, config.MaxPeers),
		maxPeers:              config.MaxPeers,
		maxConcurrencyPerPeer: config.MaxConcurrencyPerPeer,
	}
}

func (t *table) ReceivePeers(peers ...krpc.NodeAddr) (received uint, discarded uint) {
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
		if existingPeer, ok := t.peers[nodeAddr.String()]; !ok {
			peersToAdd = append(peersToAdd, nodeAddr)
		} else {
			existingPeer.lastDiscoveredAt = now
			existingPeer.discoveredCount++
		}
	}
	nPeersToAdd := len(peersToAdd)
	maxPeersToAdd := int(t.maxPeers) - len(t.peers)
	if nPeersToAdd > maxPeersToAdd {
		nEvicted := t.tryEvictPeersLocked(nPeersToAdd - maxPeersToAdd)
		nPeersToAdd = maxPeersToAdd + nEvicted
		discarded = uint(nPeersToAdd - nEvicted)
	}
	for i := 0; i < nPeersToAdd; i++ {
		nodeAddr := peersToAdd[i]
		newPeer := t.newPeer(nodeAddr)
		t.peers[nodeAddr.String()] = newPeer
		received++
	}
	return
}

func (t *table) newPeer(nodeAddr krpc.NodeAddr) *peer {
	now := time.Now()
	return &peer{
		peerSemaphore:    semaphore.NewWeighted(int64(t.maxConcurrencyPerPeer)),
		tableSemaphore:   t.semaphore,
		nodeAddr:         nodeAddr,
		discoveredAt:     now,
		lastDiscoveredAt: now,
		discoveredCount:  1,
	}
}

// tryEvictPeersLocked tries to evict n peers from the routing table based on the oldest peers to receive a successful response that can be unlocked.
// Will only consider peers with at least one response or error.
// Assumes that the caller has already locked the routing table.
// todo: Improve eviction strategy.
func (t *table) tryEvictPeersLocked(n int) int {
	maxCandidates := n * 5
	if maxCandidates > len(t.peers) {
		maxCandidates = len(t.peers)
	}
	candidates := make([]*peer, 0, maxCandidates)
	for _, peer := range t.peers {
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
			delete(t.peers, candidates[i].nodeAddr.String())
			nEvicted++
		}
		candidates[i].peerSemaphore.Release(int64(t.maxConcurrencyPerPeer))
	}
	return nEvicted
}

func (t *table) TryEachPeer(
	ctx context.Context,
	fn func(ctx context.Context, peer PeerInfo) error,
) error {
	t.mutex.RLock()
	peersToTry := make([]*peer, 0, len(t.peers))
	for _, peer := range t.peers {
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
		err = fmt.Errorf("no peer available (total %d)", len(peersToTry))
	}
	return err
}

func (t *table) tryPeerLocked(
	ctx context.Context,
	peer *peer,
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
	thePeer, ok := t.peers[info.String()]
	if !ok {
		diff := len(t.peers) - int(t.maxPeers)
		if diff > 0 {
			_ = t.tryEvictPeersLocked(diff)
		}
		thePeer = t.newPeer(info)
		t.peers[info.String()] = thePeer
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

type peer struct {
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

func (p *peer) Addr() krpc.NodeAddr {
	return p.nodeAddr
}

func (p *peer) ResponseCount() uint {
	return p.responseCount
}

func (p *peer) LastRespondedAt() time.Time {
	return p.lastRespondedAt
}

func (p *peer) WithLock(ctx context.Context, fn func(ctx context.Context) error) error {
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

func (p *peer) tryWithLock(ctx context.Context, fn func(ctx context.Context) error) (bool, error) {
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

func (p *peer) doLocked(ctx context.Context, fn func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if p.evicted {
		return errors.New("peer evicted")
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
