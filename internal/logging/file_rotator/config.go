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
		param.Dynamic(
			param.Description[Level]("The logging level"),
			param.EnumValues(slice.Map(level.LevelValues(), func(lvl level.Level) Level {
				return Level(lvl)
			})...),
			param.Default(Level(level.LevelInfo)),
		),
	)

	ParamSubPath = param.MustNew(
		param.Description[SubPath]("Subpath within the data folder to write log files"),
		param.Default(SubPath("logs")),
		param.Required[SubPath](),
	)

	ParamBaseName = param.MustNew(
		param.Description[BaseName]("Base name for log files"),
		param.Default(BaseName("bitmagnet")),
		param.Required[BaseName](),
	)

	ParamMaxAge = param.MustNew(
		param.Description[MaxAge]("Maximum age of log files to retain"),
		param.Duration[MaxAge](true),
		param.Default(MaxAge(time.Minute*60)),
	)

	ParamMaxSize = param.MustNew(
		param.Description[MaxSize]("Maximum size of a log file before writing to a new file"),
		param.Int[MaxSize](),
		param.Default(MaxSize(100_000_000)),
		param.GreaterThan(MaxSize(0)),
	)

	ParamMaxBackups = param.MustNew(
		param.Description[MaxBackups]("Maximum number of log files to retain before deleting"),
		param.Int[MaxBackups](),
		param.Default(MaxBackups(5)),
		param.Min(MaxBackups(0)),
	)

	ParamBufferSize = param.MustNew(
		param.Description[BufferSize]("Maximum number of logs to keep in memory before writing"),
		param.Int[BufferSize](),
		param.Default(BufferSize(1_000)),
		param.GreaterThan(BufferSize(0)),
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
