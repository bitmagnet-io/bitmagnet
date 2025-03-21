package torznab

import "strings"

type Profile struct {
  Name                    string `validate:"required"`
  Title                   string
  DisableOrderByRelevance bool
  DefaultLimit            uint
  MaxLimit                uint
  Tags                    []string
}

func (p Profile) MergeDefaults() Profile {
  if p.Title == "" {
    p.Title = ProfileDefault.Title
  }
  if p.DefaultLimit == 0 {
    p.DefaultLimit = ProfileDefault.DefaultLimit
  }
  if p.MaxLimit == 0 {
    p.MaxLimit = ProfileDefault.MaxLimit
  }
  return p
}

type Config struct {
  Profiles []Profile
}

var ProfileDefault = Profile{
  Name:         "default",
  Title:        "bitmagnet",
  DefaultLimit: 100,
  MaxLimit:     100,
}

func NewDefaultConfig() Config {
  return Config{}
}

func (c *Config) GetProfile(name string) (Profile, bool) {
  for _, p := range c.Profiles {
    if strings.EqualFold(p.Name, name) {
      return p, true
    }
  }
  return Profile{}, false
}
