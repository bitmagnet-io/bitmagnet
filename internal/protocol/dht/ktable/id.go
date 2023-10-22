package ktable

import "github.com/bitmagnet-io/bitmagnet/internal/protocol"

type ID = protocol.ID

func bucketIndex(a ID, b ID) int {
	return bucketIndexOffset(a, b, 0)
}

func bucketIndexOffset(a, b ID, offset int) int {
	return a.Distance(b).LeadingZeros(offset)
}
