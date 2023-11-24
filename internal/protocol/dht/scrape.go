package dht

import (
	"crypto/sha1"
	"encoding/binary"
	"github.com/bits-and-blooms/bloom/v3"
	"math"
	"math/bits"
	"net"
)

const (
	m = 256 * 8
	k = 2
)

type ScrapeBloomFilter [256]byte

// Note that if you intend for an IP to be in the IPv4 space, you might want to trim it to 4 bytes
// with IP.To4.
func (me *ScrapeBloomFilter) AddIp(ip net.IP) {
	h := sha1.New()
	h.Write(ip)
	var sum [20]byte
	h.Sum(sum[:0])
	me.addK(int(sum[0]) | int(sum[1])<<8)
	me.addK(int(sum[2]) | int(sum[3])<<8)
}

func (me *ScrapeBloomFilter) addK(index int) {
	index %= m
	me[index/8] |= 1 << (index % 8)
}

func (me ScrapeBloomFilter) countZeroes() (ret int) {
	for _, i := range me {
		ret += 8 - bits.OnesCount8(i)
	}
	return
}

func (me *ScrapeBloomFilter) EstimateCount() float64 {
	if me == nil {
		return 0
	}
	c := float64(me.countZeroes())
	if c > m-1 {
		c = m - 1
	}
	return math.Log(c/m) / (k * math.Log(1.-1./m))
}

const (
	size     = 32
	byteSize = size * 8
	M        = byteSize * 8
	K        = 2
)

func (me *ScrapeBloomFilter) ToBloomFilter() *bloom.BloomFilter {
	return bloom.FromWithM(convertBytes(*me), M, K)
}

func convertBytes(b [byteSize]byte) []uint64 {
	ret := make([]uint64, size)
	for i := 0; i < size; i++ {
		startPos := i * 8
		ret[i] = binary.BigEndian.Uint64(b[startPos : startPos+8])
	}
	return ret
}
