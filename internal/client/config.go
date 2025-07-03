package client

import (
	clientmodel "github.com/bitmagnet-io/bitmagnet/internal/client/model"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type SendTo struct {
	ID       clientmodel.ID
	Host     string
	Port     string
	Username string
	Password string
}

type Config struct {
	Enabled         bool
	SendTo          []SendTo
	DefaultCategory string
	Categories      map[model.ContentType]string
}

func NewDefaultConfig() Config {
	cfg := Config{
		Enabled:         false,
		DefaultCategory: "prowlarr",
	}
	cat := make(map[model.ContentType]string)
	cat[model.ContentTypeTvShow] = "sonarr"
	cat[model.ContentTypeMovie] = "radarr"
	cfg.Categories = cat

	return cfg
}

func (c Config) GetSendTo(id clientmodel.ID) (SendTo, bool) {
	for _, c := range c.SendTo {
		if c.ID == id {
			return c, true
		}
	}

	return SendTo{}, false
}

func (c Config) All() []clientmodel.ID {
	all := make([]clientmodel.ID, 0)

	for _, s := range c.SendTo {
		if s.ID.IsValid() {
			all = append(all, s.ID)
		}
	}

	return all
}
