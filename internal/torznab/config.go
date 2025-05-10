package torznab

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Config struct {
	Profiles []Profile
}

func (c Config) MergeDefaults() Config {
	c.Profiles = slice.Map(c.Profiles, func(profile Profile) Profile {
		return profile.MergeDefaults()
	})

	return c
}

func NewDefaultConfig() Config {
	return Config{}
}

func (c Config) GetProfile(name string) (Profile, bool) {
	for _, p := range c.Profiles {
		if strings.EqualFold(p.ID, name) {
			return p, true
		}
	}

	return Profile{}, false
}
