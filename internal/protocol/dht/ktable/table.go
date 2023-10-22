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
	PeerID ID `name:"peer_id"`
}

type Result struct {
	fx.Out
	Table TableBatch
}

func New(p Params) Result {
	return Result{
		Table: newBatcher(newTable(p.PeerID, 160, 160)),
	}
}

type TableOrigin interface {
	Origin() ID
}

type TableCommand interface {
	PutPeer(id ID, addr netip.AddrPort, options ...PeerOption) btree.PutResult
	DropPeer(id ID, reason error) bool
	PutHash(id ID, peers []HashPeer, options ...HashOption) btree.PutResult
}

type TableQuery interface {
	GetClosestPeers(id ID) []Peer
	GetOldestPeers(cutoff time.Time, n int) []Peer
	GetPeersForSampleInfoHashes(n int) []Peer
	FilterKnownAddrs(addrs []netip.Addr) []netip.Addr
	// GeneratePeerID generates a pseudo-random peer ID, biased towards the emptiest peer buckets.
	GeneratePeerID() ID
	GetHashOrClosestPeers(id ID) GetHashOrClosestPeersResult
	// SampleHashesAndPeers returns a random sample of up to 8 hashes and peers, and the total hashes count.
	SampleHashesAndPeers() SampleHashesAndPeersResult
}

type TableBatch interface {
	TableOrigin
	TableQuery
	BatchCommand(commands ...Command)
	Stats() Stats
}

type Table interface {
	TableOrigin
	TableCommand
	TableQuery
	Stats() Stats
}

type GetHashOrClosestPeersResult struct {
	Hash         Hash
	ClosestPeers []Peer
	Found        bool
}

type SampleHashesAndPeersResult struct {
	Hashes      []Hash
	Peers       []Peer
	TotalHashes int
}

func newTable(origin ID, peersK int, hashesK int) *table {
	rm := &reverseMap{addrs: make(map[string]*infoForAddr)}
	return &table{
		origin:  origin,
		peersK:  peersK,
		hashesK: hashesK,
		peers: peerBucketRoot{
			bucketRoot: newBucketRoot[netip.AddrPort, PeerOption, Peer, *peer](
				origin,
				peersK,
				func(id ID, addr netip.AddrPort) *peer {
					return &peer{
						id:           id,
						addr:         addr,
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
	peers   peerBucketRoot
	hashes  hashBucketRoot
	addrs   *reverseMap
}

func (t *table) Origin() ID {
	return t.origin
}
