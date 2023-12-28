package model

import (
	"errors"
	"fmt"
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
	tsvParts := NewTsvParts()
	tsvParts.Add(TsvPartLabelA, tc.Torrent.Name)
	tsvParts.Add(TsvPartLabelD, tc.InfoHash.String())
	if !tc.ContentID.Valid {
		titleParts = append(titleParts, tc.Torrent.Name)
	} else {
		titleParts = append(titleParts, tc.Content.Title)
		tsvParts.Add(TsvPartLabelA, tc.Content.Title)
		tc.ContentType = NewNullContentType(tc.Content.Type)
		tc.ContentSource = NewNullString(tc.Content.Source)
		tc.ContentID = NewNullString(tc.Content.ID)
		if tc.Content.OriginalTitle.Valid && tc.Content.Title != tc.Content.OriginalTitle.String {
			titleParts = append(titleParts, fmt.Sprintf("/ %s", tc.Content.OriginalTitle.String))
			tsvParts.Add(TsvPartLabelA, tc.Content.OriginalTitle.String)
		}
		tc.ReleaseDate = tc.Content.ReleaseDate
		tc.ReleaseYear = tc.Content.ReleaseYear
		for _, c := range tc.Content.Collections {
			if c.Type == "genre" {
				tsvParts.Add(TsvPartLabelD, c.Name)
			}
		}
		for _, a := range tc.Content.Attributes {
			if a.Key == "id" {
				tsvParts.Add(TsvPartLabelD, a.Value)
			}
		}
	}
	if !tc.ReleaseYear.IsNil() {
		titleParts = append(titleParts, fmt.Sprintf("(%d)", tc.ReleaseYear))
		tsvParts.Add(TsvPartLabelB, tc.ReleaseYear.String())
	}
	if tc.VideoResolution.Valid {
		tsvParts.Add(TsvPartLabelC, tc.VideoResolution.VideoResolution.Label())
	}
	if tc.VideoSource.Valid {
		tsvParts.Add(TsvPartLabelC, tc.VideoSource.VideoSource.String())
	}
	if tc.VideoCodec.Valid {
		tsvParts.Add(TsvPartLabelC, tc.VideoCodec.VideoCodec.String())
	}
	if tc.Video3d.Valid {
		tsvParts.Add(TsvPartLabelC, "3D")
	}
	if tc.VideoModifier.Valid {
		tsvParts.Add(TsvPartLabelC, tc.VideoModifier.VideoModifier.String())
	}
	if tc.ReleaseGroup.Valid {
		tsvParts.Add(TsvPartLabelC, tc.ReleaseGroup.String)
	}
	if len(tc.Languages) == 0 && tc.Content.OriginalLanguage.Valid {
		tc.Languages = Languages{tc.Content.OriginalLanguage.Language: struct{}{}}
	}
	if len(tc.Episodes) > 0 {
		titleParts = append(titleParts, tc.Episodes.String())
	}
	tc.Title = strings.Join(titleParts, " ")
	tc.TsvParts = tsvParts
	return nil
}
