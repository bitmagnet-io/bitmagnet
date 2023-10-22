package ktable

import "github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"

type Stats struct {
	//origin              string
	//PeersK              int
	//MinPeers            int
	//MaxPeers            int
	//NonEmptyPeerBuckets int
	//GrownPeerBuckets    int
	//MaxPeersK           int
	//MedianPeers int
	//Quartiles   stats.Quartiles
	TotalPeers int
	//HashesK     int
	TotalHashes int
	btree.Stats
}

func (t *table) Stats() Stats {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return Stats{
		TotalPeers:  t.peers.count(),
		TotalHashes: t.hashes.count(),
		Stats:       t.peers.btree.Stats(),
	}
}
