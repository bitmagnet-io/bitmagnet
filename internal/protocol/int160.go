package protocol

import (
	"encoding/hex"
	"math"
	"math/big"
)

type Int160 struct {
	bits [20]uint8
}

func (i Int160) String() string {
	return hex.EncodeToString(i.bits[:])
}

func (i Int160) AsByteArray() [20]byte {
	return i.bits
}

func (i Int160) ByteString() string {
	return string(i.bits[:])
}

func (i Int160) BitLen() int {
	var a big.Int

	a.SetBytes(i.bits[:])

	return a.BitLen()
}

// func (me *Int160) SetBytes(b []byte) {
//	nBuckets := copy(me.bits[:], b)
//	if nBuckets != 20 {
//		panic(nBuckets)
//	}
//}

func (i Int160) WithBit(index int, val bool) Int160 {
	var orVal uint8
	if val {
		orVal = 1 << (7 - index%8)
	}

	var mask uint8 = ^(1 << (7 - index%8))
	i.bits[index/8] = i.bits[index/8]&mask | orVal

	return i
}

func (i Int160) GetBit(index int) bool {
	return i.bits[index/8]>>(7-index%8)&1 == 1
}

func (i Int160) Bytes() []byte {
	return i.bits[:]
}

func (i Int160) Cmp(r Int160) int {
	for b := range i.bits {
		if i.bits[b] < r.bits[b] {
			return -1
		} else if i.bits[b] > r.bits[b] {
			return 1
		}
	}

	return 0
}

func (i Int160) WithMax() Int160 {
	for b := range i.bits {
		i.bits[b] = math.MaxUint8
	}

	return i
}

func (i Int160) Xor(b1, b2 Int160) Int160 {
	for b := range i.bits {
		i.bits[b] = b1.bits[b] ^ b2.bits[b]
	}

	return i
}

func (i Int160) IsZero() bool {
	for _, b := range i.bits {
		if b != 0 {
			return false
		}
	}

	return true
}

func NewInt160FromByteArray(b [20]byte) (ret Int160) {
	copy(ret.bits[:], b[:])
	return
}

func (i Int160) Distance(b Int160) (ret Int160) {
	return ret.Xor(i, b)
}
