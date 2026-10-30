package model

import (
	"strings"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pgMaxTsvectorBytes = 1048575

// pgLexemeBytes is the space PostgreSQL's tsvector input uses for a lexeme with these positions
// (uniqueentry in tsvector.c): positions are clamped to 16383, deduplicated and capped at 256;
// a lexeme with positions is aligned to 2 bytes and followed by the count and the positions.
func pgLexemeBytes(lexeme string, positions []int) int {
	if len(positions) == 0 {
		return len(lexeme)
	}

	unique := make(map[int]struct{})
	for _, p := range positions {
		unique[min(p, 16383)] = struct{}{}
	}

	return (len(lexeme)+1)/2*2 + 2 + 2*min(len(unique), 256)
}

func pgTsvectorBytes(tsv fts.Tsvector) int {
	n := 0

	for lexeme, labels := range tsv {
		positions := make([]int, 0, len(labels))
		for p := range labels {
			positions = append(positions, p)
		}

		n += pgLexemeBytes(lexeme, positions)
	}

	return n
}

// distinctWord returns a distinct 60-letter word for i.
func distinctWord(i int) string {
	word := []byte(strings.Repeat("x", 60))

	for j := 52; j < len(word); j++ {
		word[j] = byte('a' + i%26)
		i /= 26
	}

	return string(word)
}

func TestUpdateTsvLimitsHugeFileLists(t *testing.T) {
	t.Parallel()

	files := make([]TorrentFile, 20000)
	positions := make(map[string][]int)
	pos := 0

	for i := range files {
		files[i] = TorrentFile{Path: "dir/" + distinctWord(i) + ".bin"}

		for _, lexeme := range fts.TokenizeFlat(files[i].Path) {
			pos++
			positions[lexeme] = append(positions[lexeme], pos)
		}
	}

	uncapped := 0
	for lexeme, p := range positions {
		uncapped += pgLexemeBytes(lexeme, p)
	}

	require.Greater(t, uncapped, pgMaxTsvectorBytes, "the test files must exceed PostgreSQL's limit")

	tc := TorrentContent{Torrent: Torrent{Name: "pillars of eternity", Files: files}}
	tc.UpdateTsv()

	n := pgTsvectorBytes(tc.Tsv)
	assert.LessOrEqual(t, n, pgMaxTsvectorBytes)
	assert.Greater(t, n, 900_000)

	for _, lexeme := range []string{"pillars", "eternity", "dir", distinctWord(0)} {
		assert.Contains(t, tc.Tsv, lexeme)
	}

	assert.NotContains(t, tc.Tsv, distinctWord(len(files)-1))
}

func TestUpdateTsvKeepsSmallFileLists(t *testing.T) {
	t.Parallel()

	tc := TorrentContent{Torrent: Torrent{Name: "album", Files: []TorrentFile{
		{Path: "album/one.flac"}, {Path: "album/two.flac"}, {Path: "album/three.flac"},
	}}}
	tc.UpdateTsv()

	for _, lexeme := range []string{"album", "one", "two", "three", "flac"} {
		assert.Contains(t, tc.Tsv, lexeme)
	}
}

func TestAddFileSearchStringsMatchesAddText(t *testing.T) {
	t.Parallel()

	strs := []string{"Season 1/Episode 01.mkv", "", "...", "Episode 02.mkv", "sub/Episode 02.srt", "x x x"}

	want := fts.Tsvector{}
	want.AddText("some infohash", fts.TsvectorWeightA)

	got := addFileSearchStrings(want.Copy(), strs)

	for _, str := range strs {
		want.AddText(str, fts.TsvectorWeightD)
	}

	assert.Equal(t, want, got)
}
