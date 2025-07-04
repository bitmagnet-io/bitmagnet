package filerotator

import (
	"path"
	"time"

	"github.com/adrg/xdg"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
)

type Config struct {
	Enabled    bool
	Level      level.Level
	Path       string
	BaseName   string
	MaxAge     time.Duration
	MaxSize    int
	MaxBackups int
	BufferSize int
}

func NewDefaultConfig() Config {
	return Config{
		Enabled:    false,
		Level:      level.LevelInfo,
		Path:       path.Join(xdg.DataHome, "bitmagnet", "logs"),
		BaseName:   "bitmagnet",
		MaxAge:     time.Minute * 60,
		MaxSize:    1_000_000 * 100,
		BufferSize: 1_000,
		MaxBackups: 5,
	}
}
