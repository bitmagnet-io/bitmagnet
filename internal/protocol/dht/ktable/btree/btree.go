package btree

import (
	"encoding/hex"
	"errors"
)

// Btree is the binary tree implementation used by the Kademlia routing table.
// The Kademlia literature refers to "buckets" that can be "split" according to certain rules when they are full.
// This implementation does not use buckets exactly as described in the literature, but rather a simpler binary tree,
// however the end result is largely equivalent.
type Btree interface {
	N() int
	Put(NodeID) PutResult
	Has(NodeID) bool
	Drop(NodeID) bool
	Closest(NodeID, int) []NodeID
	Count() int
}

type NodeID []byte

type Bit bool

const (
	Bit0 Bit = false
	Bit1 Bit = true
)

type Bits []Bit

func (b Bits) LeadingZeros() int {
	for i, bit := range b {
		if bit {
			return i
		}
	}
	return len(b)
}

func (b Bits) Cmp(other Bits) int {
	thisPos := b.LeadingZeros() + 1
	otherPos := other.LeadingZeros() + 1
	if thisPos < otherPos {
		return 1
	}
	if thisPos > otherPos {
		return -1
	}
	if thisPos == len(b) {
		return 0
	}
	return b[thisPos:].Cmp(other[thisPos:])
}

func (b Bits) String() string {
	str := ""
	for _, bit := range b {
		if bit {
			str += "1"
		} else {
			str += "0"
		}
	}
	return str
}

func ParseBinaryNodeID(str string) (NodeID, error) {
	if len(str)%8 != 0 {
		return nil, errors.New("length must be multiple of 8")
	}
	id := make(NodeID, len(str)/8)
	for i := 0; i < len(str); i++ {
		if str[i] == '1' {
			id[i/8] |= 1 << (7 - uint(i%8))
		} else if str[i] != '0' {
			return nil, errors.New("invalid character")
		}
	}
	return id, nil
}

func MustParseBinaryNodeID(str string) NodeID {
	id, err := ParseBinaryNodeID(str)
	if err != nil {
		panic(err)
	}
	return id
}

func (id NodeID) GetBit(n int) Bit {
	return id[n/8]>>(7-uint(n%8))&1 == 1
}

func (id NodeID) Xor(other NodeID) (NodeID, error) {
	if len(id) != len(other) {
		return nil, errors.New("length mismatch")
	}
	ret := make(NodeID, len(id))
	for i := 0; i < len(id); i++ {
		ret[i] = id[i] ^ other[i]
	}
	return ret, nil
}

func (id NodeID) MustXor(other NodeID) NodeID {
	ret, err := id.Xor(other)
	if err != nil {
		panic(err)
	}
	return ret
}

func (id NodeID) String() string {
	return hex.EncodeToString(id)
}

func (id NodeID) BinaryString() string {
	return id.Bits().String()
}

func (id NodeID) Equals(other NodeID) bool {
	if len(other) != len(id) {
		return false
	}
	for i, b := range other {
		if b != id[i] {
			return false
		}
	}
	return true
}

func (id NodeID) Bits() Bits {
	path := make(Bits, len(id)*8)
	for i := 0; i < len(id)*8; i++ {
		if id.GetBit(i) {
			path[i] = Bit1
		}
	}
	return path
}

type PutResult int

const (
	PutRejected PutResult = iota
	PutAccepted
	PutAlreadyExists
)

func (r PutResult) String() string {
	switch r {
	case PutRejected:
		return "rejected"
	case PutAccepted:
		return "accepted"
	case PutAlreadyExists:
		return "already exists"
	default:
		return "unknown"
	}
}

func appendToPath(path []Bit, bit Bit) []Bit {
	newPath := make([]Bit, 0, len(path)+1)
	newPath = append(newPath, path...)
	newPath = append(newPath, bit)
	return newPath
}
