package torznab

type Profile struct {
	Name                    string `validate:"required"`
	DisableOrderByRelevance bool
	LogRequest              bool
	Tags                    []string
}

type Config struct {
	DefaultProfile Profile
	Profiles       []Profile
}

func NewDefaultConfig() Config {
	return Config{
		DefaultProfile: Profile{Name: "default"},
	}

}

func (c *Config) Map() map[string]Profile {
	profileMap := make(map[string]Profile, len(c.Profiles))
	for _, profile := range c.Profiles {
		profileMap[profile.Name] = profile
	}

	return profileMap
}
