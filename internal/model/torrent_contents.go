package model

import (
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
)

const maxPostgresTsvectorBytes = 1<<20 - 1

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

func (tc TorrentContent) baseTsv() fts.Tsvector {
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

	return tsv
}

func buildTruncatedTsvector(base fts.Tsvector, fileSearchStrings []string, maxBytes int) fts.Tsvector {
	tsv := base.Copy()

	for _, str := range fileSearchStrings {
		next := tsv.Copy()
		next.AddText(str, fts.TsvectorWeightD)
		if len(next.String()) > maxBytes {
			break
		}

		tsv = next
	}

	return tsv
}

func (tc *TorrentContent) updateTsv(maxBytes int) {
	baseTsv := tc.baseTsv()
	fileSearchStrings := tc.Torrent.fileSearchStrings()

	tsv := baseTsv.Copy()
	for _, str := range fileSearchStrings {
		tsv.AddText(str, fts.TsvectorWeightD)
	}

	if len(tsv.String()) > maxBytes {
		tsv = buildTruncatedTsvector(baseTsv, fileSearchStrings, maxBytes)
	}

	tc.Tsv = tsv
}

func (tc *TorrentContent) UpdateTsv() {
	tc.updateTsv(maxPostgresTsvectorBytes)
}
