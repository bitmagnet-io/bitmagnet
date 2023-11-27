package bloom

//
//import (
//	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
//	"testing"
//)
//
//func TestNewDefaultStableBloomFilter(t *testing.T) {
//	hashes := make(map[protocol.ID]struct{})
//	bf := NewDefaultStableBloomFilter()
//	bytes, _ := bf.GobEncode()
//	t.Log(len(bytes))
//	for i := 0; i < 10_000_000; i++ {
//		id := protocol.RandomNodeID()
//		bf.Add(id[:])
//	}
//	for i := 0; i < 1_000_000; i++ {
//		id := protocol.RandomNodeID()
//		hashes[id] = struct{}{}
//		bf.Add(id[:])
//	}
//	t.Log("added")
//	fn := 0
//	for h := range hashes {
//		if !bf.Test(h[:]) {
//			fn++
//		}
//	}
//	t.Log("fn", fn)
//	fps := make(map[protocol.ID]struct{})
//	for i := 0; i < 10_000_000; i++ {
//		id := protocol.RandomNodeID()
//		if bf.Test(id[:]) {
//			fps[id] = struct{}{}
//		}
//	}
//	for i := 0; i < 10_000_000; i++ {
//		id := protocol.RandomNodeID()
//		bf.Add(id[:])
//	}
//	t.Log("fp", len(fps))
//	unFp := 0
//	for h := range fps {
//		if !bf.Test(h[:]) {
//			unFp++
//		}
//	}
//	t.Log("unFp", unFp)
//	t.Log("sp", bf.StablePoint())
//}
