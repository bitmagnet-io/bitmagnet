package dhtcrawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveBootstrapAddrNormalizesIPv4MappedIPv6(t *testing.T) {
	t.Parallel()

	addr, err := ResolveBootstrapAddr("[::ffff:192.168.1.6]:40825")
	require.NoError(t, err)

	assert.True(t, addr.Addr().Is4())
	assert.Equal(t, "192.168.1.6:40825", addr.String())
}

func TestResolveBootstrapAddrPreservesIPv6(t *testing.T) {
	t.Parallel()

	addr, err := ResolveBootstrapAddr("[2001:db8::1]:40825")
	require.NoError(t, err)

	assert.True(t, addr.Addr().Is6())
	assert.Equal(t, "[2001:db8::1]:40825", addr.String())
}
