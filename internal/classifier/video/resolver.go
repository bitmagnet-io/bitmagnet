package video

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	TmdbClient tmdb.Client
}

type Result struct {
	fx.Out
	Resolver classifier.SubResolver `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: videoResolver{
			config:     classifier.SubResolverConfig{Key: "video", Priority: 1},
			tmdbClient: p.TmdbClient,
		},
	}
}

type videoResolver struct {
	config     classifier.SubResolverConfig
	tmdbClient tmdb.Client
}

func (r videoResolver) Config() classifier.SubResolverConfig {
	return r.config
}

func (r videoResolver) PreEnrich(content model.TorrentContent) (model.TorrentContent, error) {
	return PreEnrich(content)
}

func (r videoResolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	if content.ContentType.Valid {
		switch content.ContentType.ContentType {
		case model.ContentTypeMovie:
			return r.resolveMovie(ctx, content)
		case model.ContentTypeTvShow:
			return r.resolveTvShow(ctx, content)
		}
	}
	return model.TorrentContent{}, classifier.ErrNoMatch
}

func (r videoResolver) resolveMovie(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	externalIds := content.ExternalIds.OrderedEntries()
	if len(externalIds) > 0 {
		for _, id := range externalIds {
			if movie, err := r.tmdbClient.GetMovieByExternalId(ctx, id.Key, id.Value); err == nil {
				content.Content = movie
				if err := content.UpdateFields(); err != nil {
					return model.TorrentContent{}, err
				}
				return content, nil
			} else if !errors.Is(err, tmdb.ErrNotFound) {
				return model.TorrentContent{}, err
			}
		}
	} else if !content.ReleaseYear.IsNil() {
		if movie, err := r.tmdbClient.SearchMovie(ctx, tmdb.SearchMovieParams{
			Title:                content.Title,
			Year:                 content.ReleaseYear,
			IncludeAdult:         true,
			LevenshteinThreshold: 5,
		}); err == nil {
			content.Content = movie
			if err := content.UpdateFields(); err != nil {
				return model.TorrentContent{}, err
			}
			return content, nil
		} else if !errors.Is(err, tmdb.ErrNotFound) {
			return model.TorrentContent{}, err
		}
	}
	return model.TorrentContent{}, classifier.ErrNoMatch
}

func (r videoResolver) resolveTvShow(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	externalIds := content.ExternalIds.OrderedEntries()
	if len(externalIds) > 0 {
		for _, id := range externalIds {
			if tvShow, err := r.tmdbClient.GetTvShowByExternalId(ctx, id.Key, id.Value); err == nil {
				content.Content = tvShow
				if err := content.UpdateFields(); err != nil {
					return model.TorrentContent{}, err
				}
				return content, nil
			} else if !errors.Is(err, tmdb.ErrNotFound) {
				return model.TorrentContent{}, err
			}
		}
	} else {
		if tvShow, err := r.tmdbClient.SearchTvShow(ctx, tmdb.SearchTvShowParams{
			Name:                 content.Title,
			FirstAirDateYear:     content.ReleaseYear,
			IncludeAdult:         true,
			LevenshteinThreshold: 5,
		}); err == nil {
			content.Content = tvShow
			if err := content.UpdateFields(); err != nil {
				return model.TorrentContent{}, err
			}
			return content, nil
		} else if !errors.Is(err, tmdb.ErrNotFound) {
			return model.TorrentContent{}, err
		}
	}
	return model.TorrentContent{}, classifier.ErrNoMatch
}
