package ktable

import "github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"

type Stats struct {
	TotalPeers  int
	TotalHashes int
	btree.Stats
}

func (t *table) Stats() Stats {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.stats()
}

func (t *table) stats() Stats {
	return Stats{
		TotalPeers:  t.peers.count(),
		TotalHashes: t.hashes.count(),
		Stats:       t.peers.btree.Stats(),
	}
}
