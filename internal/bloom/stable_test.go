package bloom

//
// This is a very basic benchmark for false positive/negative rates of the stable bloom filter;
// leaving it in placs as it's handy to have.
//
// stable_test.go:23: 25000091
// stable_test.go:39: fn 17355
// stable_test.go:47: fp 8865
//
// This output demonstrates:
// - The stable bloom filter can be encoded to 25MB
// - The false negative rate averages 1.7437% for the last million deleted torrents
// - At a stable state the false positive rate is 0.09023%
//
//
// import (
//	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
//	"testing"
// )
//
// func TestNewDefaultStableBloomFilter(t *testing.T) {
// 	hashes := make(map[protocol.ID]struct{})
// 	bf := NewDefaultStableBloomFilter()
// 	bytes, _ := bf.GobEncode()
// 	t.Log(len(bytes))
// 	for i := 0; i < 10_000_000; i++ {
// 		id := protocol.RandomNodeID()
// 		bf.Add(id[:])
// 	}
// 	for i := 0; i < 1_000_000; i++ {
// 		id := protocol.RandomNodeID()
// 		hashes[id] = struct{}{}
// 		bf.Add(id[:])
// 	}
// 	fn := 0
// 	for h := range hashes {
// 		if !bf.Test(h[:]) {
// 			fn++
// 		}
// 	}
// 	t.Log("fn", fn)
// 	fps := make(map[protocol.ID]struct{})
// 	for i := 0; i < 10_000_000; i++ {
// 		id := protocol.RandomNodeID()
// 		if bf.Test(id[:]) {
// 			fps[id] = struct{}{}
// 		}
// 	}
// 	t.Log("fp", len(fps))
// }
