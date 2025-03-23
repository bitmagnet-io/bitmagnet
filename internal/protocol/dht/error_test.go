package dht

import (
	"testing"

	"github.com/anacrolix/torrent/bencode"
	"github.com/stretchr/testify/require"
)

// https://github.com/anacrolix/torrent/issues/166
func TestUnmarshalBadError(t *testing.T) {
	t.Parallel()

	var e Error
	err := bencode.Unmarshal([]byte(`l5:helloe`), &e)
	require.Error(t, err)
}
