package file_rotator

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type (
	Level      level.Level
	SubPath    string
	BaseName   string
	MaxAge     time.Duration
	MaxSize    int
	MaxBackups int
	BufferSize int
)

var (
	ParamLevel = param.MustNew(
		param.WithDynamic(
			param.WithEnumValues(slice.Map(level.LevelValues(), func(lvl level.Level) Level {
				return Level(lvl)
			})...),
			param.WithDefault(Level(level.LevelInfo)),
		),
	)

	ParamSubPath = param.MustNew(
		param.WithDefault(SubPath("logs")),
		param.WithRequired[SubPath](),
	)

	ParamBaseName = param.MustNew(
		param.WithDefault(BaseName("bitmagnet")),
		param.WithRequired[BaseName](),
	)

	ParamMaxAge = param.MustNew(
		param.WithDefault(MaxAge(time.Minute*60)),
		param.WithGreaterThan(MaxAge(0)),
	)

	ParamMaxSize = param.MustNew(
		param.WithDefault(MaxSize(100_000_000)),
		param.WithGreaterThan(MaxSize(0)),
	)

	ParamMaxBackups = param.MustNew(
		param.WithDefault(MaxBackups(5)),
		param.WithMin(MaxBackups(0)),
	)

	ParamBufferSize = param.MustNew(
		param.WithDefault(BufferSize(1_000)),
		param.WithGreaterThan(BufferSize(0)),
	)
)

// type Config struct {
// 	Level      level.Level
// 	SubPath    string
// 	BaseName   string
// 	MaxAge     time.Duration
// 	MaxSize    int
// 	MaxBackups int
// 	BufferSize int
// }

// func NewDefaultConfig() Config {
// 	return Config{
// 		Level:      level.LevelInfo,
// 		SubPath:    "logs",
// 		BaseName:   "bitmagnet",
// 		MaxAge:     time.Minute * 60,
// 		MaxSize:    100_000_000,
// 		BufferSize: 1_000,
// 		MaxBackups: 5,
// 	}
// }
