package dht

import (
	"github.com/anacrolix/torrent/bencode"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://github.com/anacrolix/torrent/issues/166
func TestUnmarshalBadError(t *testing.T) {
	var e Error
	err := bencode.Unmarshal([]byte(`l5:helloe`), &e)
	require.Error(t, err)
}
