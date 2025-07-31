package model

import (
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type ContentSourceID struct {
	Source string
	ID     string
}

type TorrentContentContentRef struct {
	Type     ContentType
	SourceID Maybe[ContentSourceID]
}

type TorrentContentRef struct {
	InfoHash   protocol.ID
	ContentRef Maybe[TorrentContentContentRef]
}

func (ref TorrentContentRef) InferID() string {
	parts := []string{
		ref.InfoHash.String(),
		"?",
		"?",
		"?",
	}

	if contentRef, ok := ref.ContentRef.ValueOK(); ok {
		parts[1] = contentRef.Type.String()
		if sourceID, ok := contentRef.SourceID.ValueOK(); ok {
			parts[2] = sourceID.Source
			parts[3] = sourceID.ID
		}
	}

	return strings.Join(parts, ":")
}

func (tc TorrentContent) Ref() TorrentContentRef {
	var contentRef Maybe[TorrentContentContentRef]

	if tc.ContentType.Valid {
		ct := tc.ContentType.ContentType
		var sourceID Maybe[ContentSourceID]
		if tc.ContentSource.Valid && tc.ContentID.Valid {
			sourceID = MaybeValid(ContentSourceID{
				Source: tc.ContentSource.String,
				ID:     tc.ContentID.String,
			})
		}
		contentRef = MaybeValid(TorrentContentContentRef{
			Type:     ct,
			SourceID: sourceID,
		})
	}

	return TorrentContentRef{
		InfoHash:   tc.InfoHash,
		ContentRef: contentRef,
	}
}

func (tc TorrentContent) InferID() string {
	return tc.Ref().InferID()
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

	for _, str := range tc.Torrent.fileSearchStrings() {
		tsv.AddText(str, fts.TsvectorWeightD)
	}

	tc.Tsv = tsv
}
