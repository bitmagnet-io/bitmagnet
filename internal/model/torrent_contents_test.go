package model

import (
	"strings"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTsvKeepsAllFileSearchStringsWhenWithinLimit(t *testing.T) {
	t.Parallel()

	tc := TorrentContent{
		InfoHash: protocol.MustParseID("0102030405060708090a0b0c0d0e0f1011121314"),
		Torrent: Torrent{
			Name: "example torrent release",
			Files: []TorrentFile{
				{Path: "dir/" + strings.Repeat("a", 40) + ".txt"},
				{Path: "dir/" + strings.Repeat("b", 40) + ".txt"},
			},
		},
	}

	tc.updateTsv(10_000)

	assert.Contains(t, tc.Tsv, strings.Repeat("a", 40))
	assert.Contains(t, tc.Tsv, strings.Repeat("b", 40))
}

func TestUpdateTsvTruncatesFileSearchStringsWhenLimitExceeded(t *testing.T) {
	t.Parallel()

	firstLexeme := strings.Repeat("a", 180)
	thirdLexeme := strings.Repeat("c", 180)

	tc := TorrentContent{
		InfoHash: protocol.MustParseID("0102030405060708090a0b0c0d0e0f1011121314"),
		Torrent: Torrent{
			Name: "example torrent release",
			Files: []TorrentFile{
				{Path: "dir/" + firstLexeme + ".txt"},
				{Path: "dir/" + strings.Repeat("b", 180) + ".txt"},
				{Path: "dir/" + thirdLexeme + ".txt"},
			},
		},
	}

	tc.updateTsv(500)

	assert.LessOrEqual(t, len(tc.Tsv.String()), 500)
	assert.Contains(t, tc.Tsv, tc.InfoHash.String())
	assert.Contains(t, tc.Tsv, "example")
	assert.Contains(t, tc.Tsv, "torrent")
	assert.Contains(t, tc.Tsv, "release")
	assert.Contains(t, tc.Tsv, firstLexeme)
	assert.NotContains(t, tc.Tsv, thirdLexeme)
}
