package ktable

import (
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"net/netip"
	"strconv"
	"sync"
	"time"
)

type batcher struct {
	mutex          sync.Mutex
	commandChannel concurrency.BatchingChannel[Command]
	counter        concurrency.AtomicCounter
	queryChannel   concurrency.BatchingChannel[batchQuery]
	resultChannels batchResultChannels
	table          *table
}

func newBatcher(t *table) TableBatch {
	b := &batcher{
		commandChannel: concurrency.NewBatchingChannel[Command](make(chan Command, 2000), 1000, time.Second/4),
		queryChannel:   concurrency.NewBatchingChannel[batchQuery](make(chan batchQuery, 2000), 100, time.Second/10),
		resultChannels: batchResultChannels{
			peers:                       make(map[string]chan []Peer),
			addrs:                       make(map[string]chan []netip.Addr),
			id:                          make(map[string]chan ID),
			getHashOrClosestPeersResult: make(map[string]chan GetHashOrClosestPeersResult),
			sampleHashesAndPeersResult:  make(map[string]chan SampleHashesAndPeersResult),
		},
		table: t,
	}
	go b.batchCommands()
	go b.batchQueries()
	return b
}

func (b *batcher) batchCommands() {
	for {
		select {
		case cs, ok := <-b.commandChannel.Out():
			if !ok {
				return
			}
			b.table.mutex.Lock()
			for _, c := range cs {
				c.exec(b.table)
			}
			b.table.mutex.Unlock()
		}
	}
}

func (b *batcher) batchQueries() {
	for {
		select {
		case qs, ok := <-b.queryChannel.Out():
			if !ok {
				return
			}
			q := newBatchQuery()
			for _, v := range qs {
				q.merge(v)
			}
			b.execQuery(*q)
		}
	}
}

func (b *batcher) Origin() ID {
	return b.table.origin
}

func (b *batcher) BatchCommand(commands ...Command) {
	for _, c := range commands {
		b.commandChannel.In() <- c
	}
}

func (b *batcher) GetClosestPeers(id ID) []Peer {
	return <-b.addPeers(GetClosestPeers{ID: id})
}

func (b *batcher) GetOldestPeers(cutoff time.Time, n int) []Peer {
	return <-b.addPeers(GetOldestPeers{Cutoff: cutoff, N: n})
}

func (b *batcher) FilterKnownAddrs(addrs []netip.Addr) []netip.Addr {
	return <-b.addAddrs(FilterKnownAddrs{Addrs: addrs})
}

func (b *batcher) GetPeersForSampleInfoHashes(n int) []Peer {
	return <-b.addPeers(GetPeersForSampleInfoHashes{N: n})
}

func (b *batcher) GeneratePeerID() ID {
	return <-b.addID(GeneratePeerID{})
}

func (b *batcher) GetHashOrClosestPeers(id ID) GetHashOrClosestPeersResult {
	return <-b.addGetHashOrClosestPeersResult(GetHashOrClosestPeers{ID: id})
}

func (b *batcher) SampleHashesAndPeers() SampleHashesAndPeersResult {
	return <-b.addSampleHashesAndPeersResult(SampleHashesAndPeers{})
}

func (b *batcher) nextKey() string {
	return strconv.Itoa(b.counter.Inc(1))
}

func (b *batcher) addPeers(q Query[[]Peer]) chan []Peer {
	key := b.nextKey()
	ch := make(chan []Peer, 1)
	b.mutex.Lock()
	b.resultChannels.peers[key] = ch
	b.mutex.Unlock()
	b.queryChannel.In() <- batchQuery{
		peers: map[string]Query[[]Peer]{
			key: q,
		},
	}
	return ch
}

func (b *batcher) addAddrs(q Query[[]netip.Addr]) chan []netip.Addr {
	key := b.nextKey()
	ch := make(chan []netip.Addr, 1)
	b.mutex.Lock()
	b.resultChannels.addrs[key] = ch
	b.mutex.Unlock()
	b.queryChannel.In() <- batchQuery{
		addrs: map[string]Query[[]netip.Addr]{
			key: q,
		},
	}
	return ch
}

func (b *batcher) addID(q Query[ID]) chan ID {
	key := b.nextKey()
	ch := make(chan ID, 1)
	b.mutex.Lock()
	b.resultChannels.id[key] = ch
	b.mutex.Unlock()
	b.queryChannel.In() <- batchQuery{
		id: map[string]Query[ID]{
			key: q,
		},
	}
	return ch
}

func (b *batcher) addGetHashOrClosestPeersResult(q Query[GetHashOrClosestPeersResult]) chan GetHashOrClosestPeersResult {
	key := b.nextKey()
	ch := make(chan GetHashOrClosestPeersResult, 1)
	b.mutex.Lock()
	b.resultChannels.getHashOrClosestPeersResult[key] = ch
	b.mutex.Unlock()
	b.queryChannel.In() <- batchQuery{
		getHashOrClosestPeersResult: map[string]Query[GetHashOrClosestPeersResult]{
			key: q,
		},
	}
	return ch
}

func (b *batcher) addSampleHashesAndPeersResult(q Query[SampleHashesAndPeersResult]) chan SampleHashesAndPeersResult {
	key := b.nextKey()
	ch := make(chan SampleHashesAndPeersResult, 1)
	b.mutex.Lock()
	b.resultChannels.sampleHashesAndPeersResult[key] = ch
	b.mutex.Unlock()
	b.queryChannel.In() <- batchQuery{
		sampleHashesAndPeersResult: map[string]Query[SampleHashesAndPeersResult]{
			key: q,
		},
	}
	return ch
}

func (b *batcher) execQuery(query batchQuery) {
	b.table.mutex.RLock()
	defer b.table.mutex.RUnlock()
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var wg sync.WaitGroup
	for k, v := range query.peers {
		if ch, ok := b.resultChannels.peers[k]; ok {
			wg.Add(1)
			go (func(k string, v Query[[]Peer], ch chan []Peer) {
				defer wg.Done()
				ch <- v.execReturn(b.table)
			})(k, v, ch)
		}
	}
	for k, v := range query.addrs {
		if ch, ok := b.resultChannels.addrs[k]; ok {
			wg.Add(1)
			go (func(k string, v Query[[]netip.Addr], ch chan []netip.Addr) {
				defer wg.Done()
				ch <- v.execReturn(b.table)
			})(k, v, ch)
		}
	}
	for k, v := range query.id {
		if ch, ok := b.resultChannels.id[k]; ok {
			wg.Add(1)
			go (func(k string, v Query[ID], ch chan ID) {
				defer wg.Done()
				ch <- v.execReturn(b.table)
			})(k, v, ch)
		}
	}
	for k, v := range query.getHashOrClosestPeersResult {
		if ch, ok := b.resultChannels.getHashOrClosestPeersResult[k]; ok {
			wg.Add(1)
			go (func(k string, v Query[GetHashOrClosestPeersResult], ch chan GetHashOrClosestPeersResult) {
				defer wg.Done()
				ch <- v.execReturn(b.table)
			})(k, v, ch)
		}
	}
	for k, v := range query.sampleHashesAndPeersResult {
		if ch, ok := b.resultChannels.sampleHashesAndPeersResult[k]; ok {
			wg.Add(1)
			go (func(k string, v Query[SampleHashesAndPeersResult], ch chan SampleHashesAndPeersResult) {
				defer wg.Done()
				ch <- v.execReturn(b.table)
			})(k, v, ch)
		}
	}
	wg.Wait()
	cleanupChannels(b.resultChannels.peers)
	cleanupChannels(b.resultChannels.addrs)
	cleanupChannels(b.resultChannels.id)
	cleanupChannels(b.resultChannels.getHashOrClosestPeersResult)
	cleanupChannels(b.resultChannels.sampleHashesAndPeersResult)
}

func cleanupChannels[T any](chans map[string]chan T) {
	for k, ch := range chans {
		close(ch)
		delete(chans, k)
	}
}

type batchQuery struct {
	peers                       map[string]Query[[]Peer]
	addrs                       map[string]Query[[]netip.Addr]
	id                          map[string]Query[ID]
	getHashOrClosestPeersResult map[string]Query[GetHashOrClosestPeersResult]
	sampleHashesAndPeersResult  map[string]Query[SampleHashesAndPeersResult]
}

func newBatchQuery() *batchQuery {
	return &batchQuery{
		peers:                       make(map[string]Query[[]Peer]),
		addrs:                       make(map[string]Query[[]netip.Addr]),
		id:                          make(map[string]Query[ID]),
		getHashOrClosestPeersResult: make(map[string]Query[GetHashOrClosestPeersResult]),
		sampleHashesAndPeersResult:  make(map[string]Query[SampleHashesAndPeersResult]),
	}
}

func (q *batchQuery) merge(m batchQuery) {
	for k, v := range m.peers {
		q.peers[k] = v
	}
	for k, v := range m.addrs {
		q.addrs[k] = v
	}
	for k, v := range m.id {
		q.id[k] = v
	}
	for k, v := range m.getHashOrClosestPeersResult {
		q.getHashOrClosestPeersResult[k] = v
	}
	for k, v := range m.sampleHashesAndPeersResult {
		q.sampleHashesAndPeersResult[k] = v
	}
}

type batchResultChannels struct {
	peers                       map[string]chan []Peer
	addrs                       map[string]chan []netip.Addr
	id                          map[string]chan ID
	getHashOrClosestPeersResult map[string]chan GetHashOrClosestPeersResult
	sampleHashesAndPeersResult  map[string]chan SampleHashesAndPeersResult
}
