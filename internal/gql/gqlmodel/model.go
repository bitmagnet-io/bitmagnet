package gqlmodel

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"time"
)

type SearchQuery struct{}

type TorrentContent struct {
	ID              string
	InfoHash        model.Hash20
	ContentType     model.NullContentType
	ContentSource   model.NullString
	ContentID       model.NullString
	Title           string
	ReleaseDate     model.Date
	ReleaseYear     model.Year
	Languages       []model.Language `json:"omitempty"`
	Episodes        *Episodes
	VideoResolution model.NullVideoResolution
	VideoSource     model.NullVideoSource
	VideoCodec      model.NullVideoCodec
	Video3d         model.NullVideo3d
	VideoModifier   model.NullVideoModifier
	ReleaseGroup    model.NullString
	SearchString    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Torrent         model.Torrent
	Content         *model.Content
}

type Episodes struct {
	Label   string
	Seasons []model.Season `json:"omitempty"`
}

func NewTorrentContentFromResultItem(item search.TorrentContentResultItem) TorrentContent {
	c := TorrentContent{
		ID:              item.ID,
		InfoHash:        item.InfoHash,
		ContentType:     item.ContentType,
		ContentSource:   item.ContentSource,
		ContentID:       item.ContentID,
		Title:           item.Title,
		ReleaseDate:     item.ReleaseDate,
		ReleaseYear:     item.ReleaseYear,
		VideoResolution: item.VideoResolution,
		VideoSource:     item.VideoSource,
		VideoCodec:      item.VideoCodec,
		Video3d:         item.Video3d,
		VideoModifier:   item.VideoModifier,
		ReleaseGroup:    item.ReleaseGroup,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
		Torrent:         item.Torrent,
	}
	if item.Content.ID != "" {
		c.Content = &item.Content
	}
	languages := item.Languages.Slice()
	if len(languages) > 0 {
		c.Languages = languages
	}
	if len(item.Episodes) > 0 {
		c.Episodes = &Episodes{
			Label:   item.Episodes.String(),
			Seasons: item.Episodes.SeasonEntries(),
		}
	}
	return c
}

type TorrentSource struct {
	Key      string
	Name     string
	ImportID model.NullString
	Seeders  model.NullUint
	Leechers model.NullUint
}

func TorrentSourcesFromTorrent(t model.Torrent) []TorrentSource {
	var sources []TorrentSource
	for _, s := range t.Sources {
		sources = append(sources, TorrentSource{
			Key:      s.Source,
			Name:     s.TorrentSource.Name,
			ImportID: s.ImportID,
			Seeders:  s.Seeders,
			Leechers: s.Leechers,
		})
	}
	return sources
}
