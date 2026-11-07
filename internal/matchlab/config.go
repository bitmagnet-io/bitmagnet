package matchlab

import (
	"strings"
)

// Config controls the external MatchLab matching API.
// MATCHLAB_URL is optional: when absent the original Bitmagnet workflow is used.
type Config struct {
	URL            string
	TimeoutSeconds int `validate:"min=1,max=120"`
	MaxConcurrent  int `validate:"min=1,max=64"`
}

func NewDefaultConfig() Config {
	return Config{TimeoutSeconds: 8, MaxConcurrent: 4}
}

func (c Config) Enabled() bool { return strings.TrimSpace(c.URL) != "" }
