package model

import (
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
)

func (tc TorrentContent) InferID() string {
	parts := make([]string, 4)
	parts[0] = tc.InfoHash.String()

	if tc.ContentType.Valid {
		parts[1] = tc.ContentType.ContentType.String()
	} else {
		parts[1] = "?"
	}

	if tc.ContentSource.Valid {
		parts[2] = tc.ContentSource.String
		parts[3] = tc.ContentID.String
	} else {
		parts[2] = "?"
		parts[3] = "?"
	}

	return strings.Join(parts, ":")
}

func (tc TorrentContent) Title() string {
	if !tc.ContentID.Valid || tc.Content.Title == "" {
		return tc.Torrent.Name
	}

	var titleParts []string

	titleParts = append(titleParts, tc.Content.Title)
	if tc.Content.OriginalTitle.Valid && tc.Content.Title != tc.Content.OriginalTitle.String {
		titleParts = append(titleParts, fmt.Sprintf("/ %s", tc.Content.OriginalTitle.String))
	}

	if !tc.Content.ReleaseYear.IsNil() {
		titleParts = append(titleParts, fmt.Sprintf("(%d)", tc.Content.ReleaseYear))
	}

	if len(tc.Episodes) > 0 {
		titleParts = append(titleParts, tc.Episodes.String())
	}

	return strings.Join(titleParts, " ")
}

func (tc TorrentContent) ContentRef() Maybe[ContentRef] {
	if tc.ContentID.Valid {
		return MaybeValid(ContentRef{
			Type:   tc.ContentType.ContentType,
			Source: tc.ContentSource.String,
			ID:     tc.ContentID.String,
		})
	}

	return Maybe[ContentRef]{}
}

func (tc *TorrentContent) UpdateTsv() {
	var tsv fts.Tsvector
	if !tc.ContentID.Valid {
		tsv = fts.Tsvector{}
	} else {
		tsv = tc.Content.Tsv.Copy()
	}

	if tc.VideoResolution.Valid {
		tsv.AddText(tc.VideoResolution.VideoResolution.Label(), fts.TsvectorWeightC)
	}

	if tc.VideoSource.Valid {
		tsv.AddText(tc.VideoSource.VideoSource.String(), fts.TsvectorWeightC)
	}

	if tc.VideoCodec.Valid {
		tsv.AddText(tc.VideoCodec.VideoCodec.String(), fts.TsvectorWeightC)
	}

	if tc.Video3D.Valid {
		tsv.AddText("3D", fts.TsvectorWeightC)
	}

	if tc.VideoModifier.Valid {
		tsv.AddText(tc.VideoModifier.VideoModifier.String(), fts.TsvectorWeightC)
	}

	if tc.ReleaseGroup.Valid {
		tsv.AddText(tc.ReleaseGroup.String, fts.TsvectorWeightC)
	}

	tsv.AddText(tc.InfoHash.String(), fts.TsvectorWeightA)
	tsv.AddText(tc.Torrent.Name, fts.TsvectorWeightA)

	tc.Tsv = addFileSearchStrings(tsv, tc.Torrent.fileSearchStrings())
}

const (
	// tsvMaxBytes keeps a tsvector below PostgreSQL's limit on its lexeme and position data
	// (1048575 bytes, "string is too long for tsvector", see #410).
	tsvMaxBytes = 1_000_000
	// tsvMaxPositions is the number of positions PostgreSQL keeps per lexeme (MAXNUMPOS).
	tsvMaxPositions = 256
	// tsvPositionBytes is the size of a stored position (WordEntryPos).
	tsvPositionBytes = 2
	// tsvPositionsHeaderBytes covers the alignment padding and the position count of a lexeme.
	tsvPositionsHeaderBytes = 3
)

// tsvLexemeBytes is an upper bound of the space PostgreSQL uses for a lexeme with n positions:
// the lexeme, alignment, the position count and the positions.
func tsvLexemeBytes(lexeme string, n int) int {
	if n == 0 {
		return len(lexeme)
	}

	return len(lexeme) + tsvPositionsHeaderBytes + tsvPositionBytes*min(n, tsvMaxPositions)
}

// addFileSearchStrings adds the file search strings in order until the next one would take the
// tsvector past tsvMaxBytes; the rest of a huge file list is not indexed. Positions are the same
// as with Tsvector.AddText, which rescans all positions on each call (quadratic for huge lists).
func addFileSearchStrings(tsv fts.Tsvector, strs []string) fts.Tsvector {
	size := 0
	nextPos := 1

	for lexeme, positions := range tsv {
		size += tsvLexemeBytes(lexeme, len(positions))

		for pos := range positions {
			if pos >= nextPos {
				nextPos = pos + 1
			}
		}
	}

	for _, str := range strs {
		lexemes := fts.TokenizeFlat(str)
		if len(lexemes) == 0 {
			continue
		}

		grow := tsvGrowth(tsv, lexemes)
		if size+grow > tsvMaxBytes {
			break
		}

		size += grow

		if nextPos > 1 {
			nextPos++
		}

		for _, lexeme := range lexemes {
			if _, ok := tsv[lexeme]; !ok {
				tsv[lexeme] = make(map[int]fts.TsvectorWeight)
			}

			tsv[lexeme][nextPos] = fts.TsvectorWeightD
			nextPos++
		}
	}

	return tsv
}

// tsvGrowth is how many bytes adding these lexemes adds to tsv in PostgreSQL (an upper bound).
func tsvGrowth(tsv fts.Tsvector, lexemes []string) int {
	counts := make(map[string]int)
	for _, lexeme := range lexemes {
		counts[lexeme]++
	}

	grow := 0

	for lexeme, n := range counts {
		positions, ok := tsv[lexeme]
		before := 0

		if ok {
			before = tsvLexemeBytes(lexeme, len(positions))
		}

		grow += tsvLexemeBytes(lexeme, len(positions)+n) - before
	}

	return grow
}
