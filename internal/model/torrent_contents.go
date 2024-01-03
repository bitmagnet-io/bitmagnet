package model

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"strings"
)

func (tc TorrentContent) EntityReference() Maybe[ContentRef] {
	if tc.ContentID.Valid {
		return MaybeValid(ContentRef{
			Type:   tc.ContentType.ContentType,
			Source: tc.ContentSource.String,
			ID:     tc.ContentID.String,
		})
	}
	return Maybe[ContentRef]{}
}

func (tc *TorrentContent) UpdateFields() error {
	// check we've got access to all the associated records needed:
	if tc.ContentID.Valid && tc.EntityReference().Val != tc.Content.Ref() {
		return errors.New("invalid Content record")
	}
	if tc.Torrent.InfoHash != tc.InfoHash {
		return errors.New("missing Torrent record")
	}
	var titleParts []string
	if !tc.ContentID.Valid {
		titleParts = append(titleParts, tc.Torrent.Name)
	} else {
		titleParts = append(titleParts, tc.Content.Title)
		tc.ContentType = NewNullContentType(tc.Content.Type)
		tc.ContentSource = NewNullString(tc.Content.Source)
		tc.ContentID = NewNullString(tc.Content.ID)
		if tc.Content.OriginalTitle.Valid && tc.Content.Title != tc.Content.OriginalTitle.String {
			titleParts = append(titleParts, fmt.Sprintf("/ %s", tc.Content.OriginalTitle.String))
		}
		tc.ReleaseDate = tc.Content.ReleaseDate
		tc.ReleaseYear = tc.Content.ReleaseYear
	}
	if !tc.ReleaseYear.IsNil() {
		titleParts = append(titleParts, fmt.Sprintf("(%d)", tc.ReleaseYear))
	}
	if len(tc.Languages) == 0 && tc.Content.OriginalLanguage.Valid {
		tc.Languages = Languages{tc.Content.OriginalLanguage.Language: struct{}{}}
	}
	if len(tc.Episodes) > 0 {
		titleParts = append(titleParts, tc.Episodes.String())
	}
	tc.Title = strings.Join(titleParts, " ")
	tc.UpdateTsv()
	return nil
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
