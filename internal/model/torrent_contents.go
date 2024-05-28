package model

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"gorm.io/gorm"
	"strings"
)

func (tc *TorrentContent) AfterFind(*gorm.DB) error {
	// the following is a mitigation following the v0.9.0 release where seeders, leechers and size are sourced from the torrent_contents table
	// because these fields won't be available until after reprocessing, calculate the values here if they are missing;
	// it should be removed at some point
	if !tc.Seeders.Valid && !tc.Leechers.Valid {
		tc.Seeders = tc.Torrent.Seeders()
		tc.Leechers = tc.Torrent.Leechers()
	}
	if tc.Size == 0 {
		tc.Size = tc.Torrent.Size
	}
	return nil
}

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
	if tc.Video3d.Valid {
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
	for _, str := range tc.Torrent.fileSearchStrings() {
		tsv.AddText(str, fts.TsvectorWeightD)
	}
	tc.Tsv = tsv
}
