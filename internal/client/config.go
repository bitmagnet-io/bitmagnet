package client

import (
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type DownloadClient struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Config struct {
	Transmission    DownloadClient
	Qbittorrent     DownloadClient
	DownloadClient  gen.ClientID
	DefaultCategory string
	Categories      map[model.ContentType]string
}

func NewDefaultConfig() Config {
	cfg := Config{
		Transmission: DownloadClient{
			Host: "localhost",
			Port: "9091",
		},
		Qbittorrent: DownloadClient{
			Host:     "localhost",
			Port:     "8080",
			Username: "required",
			Password: "required",
		},
		DownloadClient:  gen.ClientIDQBittorrent,
		DefaultCategory: "prowlarr",
	}
	cat := make(map[model.ContentType]string)
	cat[model.ContentTypeTvShow] = "sonarr"
	cat[model.ContentTypeMovie] = "radarr"
	cfg.Categories = cat

	return cfg
}
