package torznab

import "strings"

type Config struct {
	Profiles []Profile
}

func (c Config) MergeDefaults() Config {
	var profiles []Profile
	for _, p := range c.Profiles {
		profiles = append(profiles, p.MergeDefaults())
	}
	c.Profiles = profiles
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
