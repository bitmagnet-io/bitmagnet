package dhtcrawler

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/bootstrap"
	"go.uber.org/fx"
)

type (
	Autostart bool

	BootstrapNodes []string

	ReseedBootstrapNodesInterval time.Duration

	SaveFilesThreshold int

	SavePieces bool

	RescrapeThreshold time.Duration
)

var (
	ParamAutostart = param.MustNew(
		param.Description[Autostart]("Start the DHT crawler automatically"),
		param.Bool[Autostart](),
		param.Default(Autostart(true)),
	)

	ParamBootstrapNodes = param.MustNew(
		param.Description[BootstrapNodes]("A list of node addresses for bootstrapping the DHT crawler"),
		param.Slice[string, BootstrapNodes](),
		param.NewDefault(func() BootstrapNodes {
			return bootstrap.DefaultBootstrapNodes
		}),
		param.MinItems[string, BootstrapNodes](1),
	)

	ParamReseedBootstrapNodesInterval = param.MustNew(
		param.Description[ReseedBootstrapNodesInterval]("Interval between reseeding of the bootstrap nodes"),
		param.Duration[ReseedBootstrapNodesInterval](true),
		param.Default(ReseedBootstrapNodesInterval(time.Minute*10)),
	)

	ParamSaveFilesThreshold = param.MustNew(
		param.Description[SaveFilesThreshold]("The maximum limit for the number of files saved by the DHT crawler with each torrent"),
		param.Int[SaveFilesThreshold](),
		param.Default(SaveFilesThreshold(100)),
		param.Min[SaveFilesThreshold](0),
	)

	ParamSavePieces = param.MustNew(
		param.Description[SavePieces]("Save torrent pieces with each torrent (uses a lot of disk space)"),
		param.Bool[SavePieces](),
	)

	ParamRescrapeThreshold = param.MustNew(
		param.Description[RescrapeThreshold]("The minimum interval before re-querying seeder and leecher numbers from the DHT"),
		param.Default(RescrapeThreshold(time.Hour*24*30)),
		param.Duration[RescrapeThreshold](true),
	)
)

type Config struct {
	fx.In
	BootstrapNodes               BootstrapNodes
	ReseedBootstrapNodesInterval ReseedBootstrapNodesInterval
	SaveFilesThreshold           SaveFilesThreshold
	SavePieces                   SavePieces
	RescrapeThreshold            RescrapeThreshold
}
