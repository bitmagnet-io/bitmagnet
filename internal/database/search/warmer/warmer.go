package warmer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
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
	Search search.Search
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	AppHook fx.Hook `group:"app_hooks"`
}

func New(params Params) Result {
	logger := params.Logger.Named("search_cache_warmer")
	stopped := make(chan struct{})
	return Result{
		AppHook: fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					ctx, cancel := context.WithCancel(context.Background())
					ticker := time.NewTicker(time.Second)
					for {
						select {
						case <-stopped:
							cancel()
							return
						case <-ticker.C:
							for _, w := range warmers.Entries() {
								logger.Infow("warming", "warmer", w.Key)
								if _, err := params.Search.TorrentContent(
									ctx,
									query.Limit(0),
									search.TorrentContentCoreJoins(),
									w.Value,
									query.CacheMode(cache.ModeWarm),
								); err != nil {
									logger.Errorw("error warming", "warmer", w.Key, "error", err)
								}
							}
							ticker.Reset(10 * time.Minute)
							continue
						}
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				close(stopped)
				return nil
			},
		},
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
			f.Value(query.FacetHasAggregationConfig(query.FacetAggregationConfig{
				TotalCount: true,
			})),
		))
	}
	for _, ct := range model.ContentTypeValues() {
		for _, f := range facets.Entries()[1:] {
			warmers.Set("aggs:"+ct.String()+"/"+f.Key, query.Options(query.WithFacet(
				search.TorrentContentTypeFacet(query.FacetHasFilter(query.FacetFilter{
					ct.String(): struct{}{},
				}))), query.WithFacet(f.Value(query.FacetHasAggregationConfig(query.FacetAggregationConfig{
				Filtered:   true,
				TotalCount: true,
			})))))
		}
	}
}
