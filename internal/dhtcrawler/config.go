package dhtcrawler

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/bootstrap"
	"go.uber.org/fx"
)

type (
	BootstrapNodes []string

	ReseedBootstrapNodesInterval time.Duration

	SaveFilesThreshold int

	SavePieces bool

	RescrapeThreshold time.Duration
)

var (
	ParamBootstrapNodes = param.MustNew(
		param.WithNewDefault(func() BootstrapNodes {
			return bootstrap.DefaultBootstrapNodes
		}),
		param.MinLength[any, string, BootstrapNodes](1),
	)

	ParamReseedBootstrapNodesInterval = param.MustNew(
		param.WithDefault(ReseedBootstrapNodesInterval(time.Minute*10)),
		param.WithMin(ReseedBootstrapNodesInterval(time.Minute)),
	)

	ParamSaveFilesThreshold = param.MustNew(
		param.WithDefault(SaveFilesThreshold(100)),
	)

	ParamSavePieces = param.MustNew[SavePieces]()

	ParamRescrapeThreshold = param.MustNew(
		param.WithDefault(RescrapeThreshold(time.Hour*24*30)),
		param.WithMin(RescrapeThreshold(time.Hour*24)),
	)
)

type Config struct {
	fx.In
	BootstrapNodes               BootstrapNodes
	ReseedBootstrapNodesInterval ReseedBootstrapNodesInterval
	// SaveFilesThreshold specifies a maximum number of files in a torrent before file information is discarded.
	// Some torrents contain thousands of files which can severely impact performance and uses a lot of disk space.
	SaveFilesThreshold SaveFilesThreshold
	// SavePieces when true, torrent pieces will be persisted to the database.
	// The pieces take up quite a lot of space, and aren't currently very useful,
	// but they may be used by future features.
	SavePieces SavePieces
	// RescrapeThreshold is the amount of time that must pass before a torrent is rescraped
	// to count seeders and leechers.
	RescrapeThreshold RescrapeThreshold
}

// func NewDefaultConfig() Config {
// 	return Config{
// 		ScalingFactor:                10,
// 		BootstrapNodes:               bootstrap.DefaultBootstrapNodes,
// 		ReseedBootstrapNodesInterval: time.Minute * 10,
// 		SaveFilesThreshold:           100,
// 		SavePieces:                   false,
// 		RescrapeThreshold:            time.Hour * 24 * 30,
// 	}
// }
