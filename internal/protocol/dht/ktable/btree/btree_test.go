package btree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testK = 4

var testOrigin = MustParseBinaryNodeID("0000111100000000")

func newTestID(str string) NodeID {
	return MustParseBinaryNodeID(str).MustXor(testOrigin)
}

// the test IDs are defined as XORs with respect to the testOrigin
var testIDs = []NodeID{
	// these should all be allowed in their respective "buckets" without any splitting:
	newTestID("0000000001001000"), // 0
	newTestID("0000000001001100"), // 1
	newTestID("0000000001001110"), // 2
	newTestID("0000000001001111"), // 3
	newTestID("0000000000110000"), // 4
	newTestID("0000000000110001"), // 5
	newTestID("0000000000110010"), // 6
	newTestID("0000000000110011"), // 7
	newTestID("1000000000001000"), // 8
	newTestID("1000000000001001"), // 9
	newTestID("1000000000001010"), // 10
	newTestID("1000000000001011"), // 11
	// these then won't be allowed with splitting disabled, but will be with splitting enabled:
	newTestID("0000000000100100"), // 12
	newTestID("0000000000100101"), // 13
	newTestID("0000000000100110"), // 14
	newTestID("0000000000100111"), // 15
	// these then won't be allowed, whether splitting is enabled or not
	newTestID("0000000000110100"), // 16
	newTestID("0000000000111000"), // 17
	newTestID("1010000000001000"), // 18
	newTestID("1010000000001001"), // 19
}

func assertPut(t *testing.T, root Btree, id NodeID, expectedResult PutResult, expectedEvicted NodeID) {
	label := "xor: " + id.MustXor(testOrigin).BinaryString()
	if expectedEvicted != nil {
		label += " / evicted: " + expectedEvicted.MustXor(testOrigin).BinaryString()
	}
	result, evicted := root.Put(id)
	assert.Equal(t, expectedResult, result, label)
	assert.Equal(t, expectedEvicted, evicted, label)
}

func TestBtree_simple(t *testing.T) {
	root := New(testOrigin, testK, false, false)
	assertPut(t, root, testOrigin, PutRejected, nil)
	for range []int{1, 2} {
		for i := 0; i < 12; i++ {
			assertPut(t, root, testIDs[i], PutAccepted, nil)
			assert.True(t, root.Has(testIDs[i]), i)
		}
		for i := 12; i < 20; i++ {
			assertPut(t, root, testIDs[i], PutRejected, nil)
		}
		for i := 0; i < 12; i++ {
			assertPut(t, root, testIDs[i], PutAlreadyExists, nil)
		}
		assert.Equal(t, 12, root.Count())
		for i := 0; i < 12; i++ {
			assert.True(t, root.Has(testIDs[i]), i)
		}
		for i := 12; i < 20; i++ {
			assert.False(t, root.Has(testIDs[i]), i)
		}
		for i := 0; i < 12; i++ {
			assert.True(t, root.Drop(testIDs[i]), i)
		}
		for i := 12; i < 20; i++ {
			assert.False(t, root.Drop(testIDs[i]), i)
		}
		assert.Equal(t, 0, root.Count())
	}
}

func TestBtree_splitting(t *testing.T) {
	root := New(testOrigin, testK, true, false)
	assertPut(t, root, testOrigin, PutRejected, nil)
	for i := 0; i < 16; i++ {
		assertPut(t, root, testIDs[i], PutAccepted, nil)
	}
	assert.Equal(t, 12, root.countCloserThan(testIDs[16]))
	for i := 16; i < 20; i++ {
		assertPut(t, root, testIDs[i], PutRejected, nil)
	}
	for i := 0; i < 16; i++ {
		assert.True(t, root.Has(testIDs[i]), i)
	}
	for i := 16; i < 20; i++ {
		assert.False(t, root.Has(testIDs[i]), i)
	}
	for i := 0; i < 16; i++ {
		assert.True(t, root.Drop(testIDs[i]), i)
	}
	for i := 16; i < 20; i++ {
		assert.False(t, root.Drop(testIDs[i]), i)
	}
	assert.Equal(t, 0, root.Count())
}

func TestBtree_eviction(t *testing.T) {
	root := New(testOrigin, testK, false, true)
	assertPut(t, root, testOrigin, PutRejected, nil)
	for i := 4; i < 8; i++ {
		assertPut(t, root, testIDs[i], PutAccepted, nil)
	}
	assertPut(t, root, testIDs[19], PutAccepted, nil)
	assertPut(t, root, testIDs[12], PutAccepted, testIDs[19])
	assertPut(t, root, testIDs[13], PutAccepted, testIDs[7])
	assertPut(t, root, testIDs[14], PutAccepted, testIDs[6])
	assertPut(t, root, testIDs[15], PutAccepted, testIDs[5])
}

func TestBtree_closest(t *testing.T) {
	root := New(testOrigin, testK, true, false)
	for i := 0; i < 16; i++ {
		assertPut(t, root, testIDs[i], PutAccepted, nil)
	}
	assert.Equal(t, []NodeID{testIDs[4], testIDs[5], testIDs[6], testIDs[7]}, root.Closest(testIDs[16], 4))
	assert.Equal(t, []NodeID{
		testIDs[8],
		testIDs[9],
		testIDs[10],
		testIDs[11],
		testIDs[12],
		testIDs[13],
		testIDs[14],
		testIDs[15],
		testIDs[4],
		testIDs[5],
		testIDs[6],
		testIDs[7],
		testIDs[0],
		testIDs[1],
		testIDs[2],
		testIDs[3],
	}, root.Closest(testIDs[18], 100))
	assert.Equal(t, []NodeID{
		testIDs[4],
		testIDs[5],
		testIDs[6],
		testIDs[7],
		testIDs[12],
		testIDs[13],
		testIDs[14],
		testIDs[15],
		testIDs[0],
		testIDs[1],
		testIDs[2],
		testIDs[3],
		testIDs[8],
		testIDs[9],
		testIDs[10],
		testIDs[11],
	}, root.Closest(testIDs[16], 100))
}
