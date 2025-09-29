package gqlmodel

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Episodes struct {
	Label   string
	Seasons []model.Season `json:"omitempty"`
}

type TorrentSourceInfo struct {
	Key      string
	Name     string
	ImportID model.NullString
	Seeders  model.NullUint
	Leechers model.NullUint
}

func TorrentSourceInfosFromTorrent(t model.Torrent) []TorrentSourceInfo {
	sources := make([]TorrentSourceInfo, 0, len(t.Sources))

	for _, s := range t.Sources {
		sources = append(sources, TorrentSourceInfo{
			Key:      s.Source,
			Name:     s.TorrentSource.Name,
			ImportID: s.ImportID,
			Seeders:  s.Seeders,
			Leechers: s.Leechers,
		})
	}

	return sources
}
