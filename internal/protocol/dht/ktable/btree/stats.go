package btree

type Stats struct {
	Buckets []BucketStats
}

type BucketStats struct {
	N     int
	Count int
}

func (n *rootNode) Stats() Stats {
	buckets := make([]BucketStats, 0, n.N())
	for i := n.N() - 1; i >= 0; i-- {
		//path := make(Bits, i+1)
		//path[i] = Bit1
		count := n.bucketCounts[i]
		if count > 0 {
			buckets = append(buckets, BucketStats{
				N:     i,
				Count: count,
			})
		}
	}
	return Stats{
		Buckets: buckets,
	}
}
