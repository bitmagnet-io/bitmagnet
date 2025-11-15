package dhtcrawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
)

func TestCreateTorrentModelSkipsPaddingFiles(t *testing.T) {
	t.Parallel()

	hash := protocol.MustParseID("0123456789abcdef0123456789abcdef01234567")
	info := metainfo.Info{
		Name: "test",
		Files: []metainfo.FileInfo{
			{
				Path:   []string{"foo.txt"},
				Length: 10,
			},
			{
				Path: []string{"attr", "padding.bin"},
				ExtendedFileAttrs: metainfo.ExtendedFileAttrs{
					Attr: "p",
				},
				Length: 20,
			},
			{
				Path:   []string{".pad", "30"},
				Length: 30,
			},
			{
				Path:   []string{"bar", "baz.bin"},
				Length: 40,
			},
		},
		PieceLength: 4,
	}

	torrent, err := createTorrentModel(hash, info, false, 10)
	require.NoError(t, err)

	require.Len(t, torrent.Files, 2)
	assert.Equal(t, uint(0), torrent.Files[0].Index)
	assert.Equal(t, uint(3), torrent.Files[1].Index)

	if assert.True(t, torrent.FilesCount.Valid) {
		assert.Equal(t, uint(2), torrent.FilesCount.Uint)
	}
	assert.Equal(t, model.FilesStatusMulti, torrent.FilesStatus)
}

func TestCreateTorrentModelCountsNonPaddingFilesPastThreshold(t *testing.T) {
	t.Parallel()

	hash := protocol.MustParseID("abcdefabcdefabcdefabcdefabcdefabcdefabcd")
	info := metainfo.Info{
		Name: "test",
		Files: []metainfo.FileInfo{
			{
				Path: []string{".pad", "1"},
				ExtendedFileAttrs: metainfo.ExtendedFileAttrs{
					Attr: "p",
				},
				Length: 1,
			},
			{
				Path:   []string{"a.bin"},
				Length: 10,
			},
			{
				Path:   []string{"b.bin"},
				Length: 20,
			},
		},
		PieceLength: 4,
	}

	torrent, err := createTorrentModel(hash, info, false, 1)
	require.NoError(t, err)

	assert.Equal(t, model.FilesStatusOverThreshold, torrent.FilesStatus)
	if assert.True(t, torrent.FilesCount.Valid) {
		assert.Equal(t, uint(2), torrent.FilesCount.Uint)
	}
	require.Len(t, torrent.Files, 1)
	assert.Equal(t, uint(1), torrent.Files[0].Index)
}
