package bloom

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestScrapeBloomFilter(t *testing.T) {
	var s dht.ScrapeBloomFilter
	assert.Equal(t, 0.5, s.EstimateCount())
	s.AddIP(net.IPv4(127, 0, 0, 1))
	c := s.EstimateCount()
	assert.Greater(t, c, 0.9)
	assert.Less(t, c, 1.1)
}

func TestConvertBloom(t *testing.T) {
	var s dht.ScrapeBloomFilter
	s.AddIP(net.IPv4(127, 0, 0, 1))
	s.AddIP(net.IPv4(127, 0, 0, 2))
	f := FromScrape(s)
	assert.Equal(t, uint32(2), f.ApproximatedSize())
}
