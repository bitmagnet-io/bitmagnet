package dhtcrawler

import (
	"testing"

	anacrolixmetainfo "github.com/anacrolix/torrent/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTorrentModelSkipsDuplicateFilePaths(t *testing.T) {
	t.Parallel()

	info := metainfo.Info{
		Name: "test",
		Files: []anacrolixmetainfo.FileInfo{
			{Path: []string{"a.txt"}, Length: 1},
			{Path: []string{"a.txt"}, Length: 2},
			{Path: []string{"b.txt"}, Length: 3},
		},
	}

	torrent, err := createTorrentModel(protocol.ID{}, info, false, 100)
	require.NoError(t, err)
	require.Len(t, torrent.Files, 2)

	assert.Equal(t, "a.txt", torrent.Files[0].Path)
	assert.Equal(t, uint(0), torrent.Files[0].Index)
	assert.Equal(t, "b.txt", torrent.Files[1].Path)
	assert.Equal(t, uint(2), torrent.Files[1].Index)
}
