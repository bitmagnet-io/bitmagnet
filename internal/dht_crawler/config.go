package dht_crawler

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
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

	PingConcurrency             int
	FindNodesConcurrency        int
	GetPeersConcurrency         int
	SampleInfoHashesConcurrency int
	ScrapeConcurrency           int
	RequestMetaInfoConcurrency  int
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

	ParamPingConcurrency = param.MustNew(
		param.Description[*atomic.Value[PingConcurrency]]("Concurrent ping workers"),
		param.Dynamic(
			param.Int[PingConcurrency](),
			param.Default(PingConcurrency(100)),
			param.Min[PingConcurrency](1),
		),
	)

	ParamFindNodesConcurrency = param.MustNew(
		param.Description[*atomic.Value[FindNodesConcurrency]]("Concurrent find_nodes workers"),
		param.Dynamic(
			param.Int[FindNodesConcurrency](),
			param.Default(FindNodesConcurrency(1000)),
			param.Min[FindNodesConcurrency](1),
		),
	)

	ParamGetPeersConcurrency = param.MustNew(
		param.Description[*atomic.Value[GetPeersConcurrency]]("Concurrent get_peers workers"),
		param.Dynamic(
			param.Int[GetPeersConcurrency](),
			param.Default(GetPeersConcurrency(1000)),
			param.Min[GetPeersConcurrency](1),
		),
	)

	ParamSampleInfoHashesConcurrency = param.MustNew(
		param.Description[*atomic.Value[SampleInfoHashesConcurrency]]("Concurrent sample_infohashes workers"),
		param.Dynamic(
			param.Int[SampleInfoHashesConcurrency](),
			param.Default(SampleInfoHashesConcurrency(1000)),
			param.Min[SampleInfoHashesConcurrency](1),
		),
	)

	ParamScrapeConcurrency = param.MustNew(
		param.Description[*atomic.Value[ScrapeConcurrency]]("Concurrent scrape workers"),
		param.Dynamic(
			param.Int[ScrapeConcurrency](),
			param.Default(ScrapeConcurrency(100)),
			param.Min[ScrapeConcurrency](1),
		),
	)

	ParamRequestMetaInfoConcurrency = param.MustNew(
		param.Description[*atomic.Value[RequestMetaInfoConcurrency]]("Concurrent meta info request workers"),
		param.Dynamic(
			param.Int[RequestMetaInfoConcurrency](),
			param.Default(RequestMetaInfoConcurrency(1000)),
			param.Min[RequestMetaInfoConcurrency](1),
		),
	)
)

type Config struct {
	fx.In
	BootstrapNodes               BootstrapNodes
	ReseedBootstrapNodesInterval ReseedBootstrapNodesInterval
	SaveFilesThreshold           SaveFilesThreshold
	SavePieces                   SavePieces
	RescrapeThreshold            RescrapeThreshold
	PingConcurrency              *atomic.Value[PingConcurrency]
	FindNodesConcurrency         *atomic.Value[FindNodesConcurrency]
	GetPeersConcurrency          *atomic.Value[GetPeersConcurrency]
	SampleInfoHashesConcurrency  *atomic.Value[SampleInfoHashesConcurrency]
	ScrapeConcurrency            *atomic.Value[ScrapeConcurrency]
	RequestMetaInfoConcurrency   *atomic.Value[RequestMetaInfoConcurrency]
}
