package warmer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/zap"
	"time"
)

type warmer struct {
	stopped  chan struct{}
	interval time.Duration
	search   search.Search
	logger   *zap.SugaredLogger
}

func (w warmer) start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			warmed := make(chan struct{})
			go func() {
				w.warm(ctx)
				close(warmed)
			}()
			// wait for warming to complete
			select {
			case <-ctx.Done():
				return
			case <-warmed:
			}
			// then wait for the next tick
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
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

func init() {

	facets := maps.NewInsertMap[string, func(options ...query.FacetOption) query.Facet]()
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

	// All the top-level facets should be warmed:
	for _, f := range facets.Entries() {
		warmers.Set("aggs:"+f.Key, query.WithFacet(
			f.Value(query.FacetIsAggregated()),
		))
	}
	// All the top-level facets within each content type should be warmed:
	for _, ct := range model.ContentTypeValues() {
		for _, f := range facets.Entries()[1:] {
			warmers.Set("aggs:"+ct.String()+"/"+f.Key, query.Options(query.WithFacet(
				search.TorrentContentTypeFacet(query.FacetHasFilter(query.FacetFilter{
					ct.String(): struct{}{},
				}))), query.WithFacet(f.Value(query.FacetIsAggregated()))))
		}
	}
}
