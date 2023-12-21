package warmer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Params struct {
	fx.In
	Config Config
	Search search.Search
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	AppHook fx.Hook `group:"app_hooks"`
}

func New(params Params) Result {
	w := warmer{
		stopped: make(chan struct{}),
		ticker:  time.NewTicker(params.Config.Interval),
		search:  params.Search,
		logger:  params.Logger.Named("search_warmer"),
	}
	return Result{
		AppHook: fx.Hook{
			OnStart: func(context.Context) error {
				go w.start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				close(w.stopped)
				return nil
			},
		},
	}
}

type warmer struct {
	stopped chan struct{}
	ticker  *time.Ticker
	search  search.Search
	logger  *zap.SugaredLogger
}

func (w warmer) start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			w.warm(ctx)
			select {
			case <-ctx.Done():
				return
			case <-w.ticker.C:
				continue
			}
		}
	}()
	<-w.stopped
}

func (w warmer) warm(ctx context.Context) {
	for _, e := range warmers.Entries() {
		w.logger.Debugw("warming", "warmer", e.Key)
		if _, err := w.search.TorrentContent(
			ctx,
			query.Limit(0),
			search.TorrentContentCoreJoins(),
			e.Value,
			query.CacheWarm(),
		); err != nil {
			w.logger.Errorw("error warming", "warmer", e.Key, "error", err)
		}
	}
}

var warmers = maps.NewInsertMap[string, query.Option]()

var facets = maps.NewInsertMap[string, func(options ...query.FacetOption) query.Facet]()

func init() {

	facets.Set(search.TorrentContentTypeFacetKey, search.TorrentContentTypeFacet)
	facets.Set(search.ContentGenreFacetKey, search.TorrentContentGenreFacet)
	facets.Set(search.LanguageFacetKey, search.TorrentContentLanguageFacet)
	facets.Set(search.Video3dFacetKey, search.Video3dFacet)
	facets.Set(search.VideoCodecFacetKey, search.VideoCodecFacet)
	facets.Set(search.VideoModifierFacetKey, search.VideoModifierFacet)
	facets.Set(search.VideoResolutionFacetKey, search.VideoResolutionFacet)
	facets.Set(search.VideoSourceFacetKey, search.VideoSourceFacet)
	facets.Set(search.TorrentFileTypeFacetKey, search.TorrentFileTypeFacet)
	facets.Set(search.TorrentSourceFacetKey, search.TorrentSourceFacet)
	facets.Set(search.TorrentTagFacetKey, search.TorrentTagsFacet)

	for _, f := range facets.Entries() {
		warmers.Set("aggs:"+f.Key, query.WithFacet(
			f.Value(query.FacetIsAggregated()),
		))
	}
	for _, ct := range model.ContentTypeValues() {
		for _, f := range facets.Entries()[1:] {
			warmers.Set("aggs:"+ct.String()+"/"+f.Key, query.Options(query.WithFacet(
				search.TorrentContentTypeFacet(query.FacetHasFilter(query.FacetFilter{
					ct.String(): struct{}{},
				}))), query.WithFacet(f.Value(query.FacetIsAggregated()))))
		}
	}
}
