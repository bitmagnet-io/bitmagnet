package dhtcrawler

import (
	"net/netip"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/stretchr/testify/require"
)

type fakeTestAndAdder struct {
	seen map[string]struct{}
}

func (f *fakeTestAndAdder) TestAndAdd(data []byte) bool {
	if f.seen == nil {
		f.seen = make(map[string]struct{})
	}

	key := string(data)
	_, ok := f.seen[key]
	f.seen[key] = struct{}{}

	return ok
}

func TestFilterDiscoveredNodesSkipsRecentlySeenAddrs(t *testing.T) {
	c := crawler{
		seenNodes: &seenNodes{
			bloom: &fakeTestAndAdder{},
		},
	}

	firstBatch := []ktable.Node{
		testNode("1111111111111111111111111111111111111111", "203.0.113.10:6881"),
		testNode("2222222222222222222222222222222222222222", "203.0.113.10:49001"),
		testNode("3333333333333333333333333333333333333333", "203.0.113.11:6881"),
	}

	filtered := c.filterDiscoveredNodes(firstBatch)
	require.Len(t, filtered, 2)
	require.Equal(t, firstBatch[0].Addr(), filtered[0].Addr())
	require.Equal(t, firstBatch[2].Addr(), filtered[1].Addr())

	secondBatch := []ktable.Node{
		testNode("4444444444444444444444444444444444444444", "203.0.113.10:6881"),
		testNode("5555555555555555555555555555555555555555", "203.0.113.12:6881"),
	}

	filtered = c.filterDiscoveredNodes(secondBatch)
	require.Len(t, filtered, 1)
	require.Equal(t, secondBatch[1].Addr(), filtered[0].Addr())
}

func TestSeenNodesNormalizesIPv4MappedIPv6(t *testing.T) {
	seen := seenNodes{
		bloom: &fakeTestAndAdder{},
	}

	require.False(t, seen.testAndAdd(netip.MustParseAddr("::ffff:203.0.113.10")))
	require.True(t, seen.testAndAdd(netip.MustParseAddr("203.0.113.10")))
}

func testNode(id string, addr string) ktable.Node {
	return ktable.NewNode(protocol.MustParseID(id), netip.MustParseAddrPort(addr))
}
