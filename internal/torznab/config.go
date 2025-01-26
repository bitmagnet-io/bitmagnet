package torznab

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
)

type Profile struct {
	Name           string `validate:"required"`
	OrderBy        search.TorrentContentOrderBy
	OrderDirection search.OrderDirection
	Tags           []string
}

type Config struct {
	DefaultProfile Profile
	Profiles       []Profile
}

// force in defaulting of []Profiles.  Name is mandated by validator
func (p Profile) completeProfile() Profile {
	orderBy := search.TorrentContentOrderByRelevance
	direction := search.OrderDirectionDescending
	if p.OrderBy != "" {
		orderBy = p.OrderBy
	}
	if p.OrderDirection != "" {
		direction = p.OrderDirection
	}
	return Profile{
		Name:           p.Name,
		OrderBy:        orderBy,
		OrderDirection: direction,
		Tags:           p.Tags,
	}
}

func NewDefaultConfig() Config {
	return Config{
		DefaultProfile: Profile{Name: "default"}.completeProfile(),
	}

}

func (c *Config) Map() map[string]Profile {
	profileMap := make(map[string]Profile, len(c.Profiles))
	for _, profile := range c.Profiles {
		profileMap[profile.Name] = profile.completeProfile()
	}

	return profileMap
}
