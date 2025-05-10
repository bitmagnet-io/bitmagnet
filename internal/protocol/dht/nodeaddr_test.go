package dht

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalNodeAddr(t *testing.T) {
	t.Parallel()

	var na NodeAddr

	require.NoError(t, na.UnmarshalBinary([]byte("\x01\x02\x03\x04\x05\x06")))
	assert.Equal(t, "1.2.3.4", na.IP.String())
}
