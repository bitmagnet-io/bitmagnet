package torznab

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Config struct {
	BaseURL  string
	Profiles []Profile
}

const configBaseURLUndefined = "__undefined__"

func (c Config) MergeDefaults() Config {
	c.Profiles = slice.Map(c.Profiles, func(profile Profile) Profile {
		return profile.MergeDefaults()
	})

	return c
}

func NewDefaultConfig() Config {
	return Config{
		BaseURL: configBaseURLUndefined,
	}
}

func (c Config) GetProfile(name string) (Profile, bool) {
	for _, p := range c.Profiles {
		if strings.EqualFold(p.ID, name) {
			if c.BaseURL != configBaseURLUndefined && !p.BaseURL.Valid {
				p.BaseURL = model.NewNullString(c.BaseURL)
			}

			return p, true
		}
	}

	return Profile{}, false
}
