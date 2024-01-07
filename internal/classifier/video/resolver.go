package video

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

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

func (r videoResolver) resolveMovie(ctx context.Context, tc model.TorrentContent) (model.TorrentContent, error) {
	externalIds := tc.ExternalIds.OrderedEntries()
	if len(externalIds) > 0 {
		for _, id := range externalIds {
			if movie, err := r.tmdbClient.GetMovieByExternalId(ctx, id.Key, id.Value); err == nil {
				if err := tc.SetContent(movie); err != nil {
					return model.TorrentContent{}, err
				}
				return tc, nil
			} else if !errors.Is(err, tmdb.ErrNotFound) {
				return model.TorrentContent{}, err
			}
		}
	} else if !tc.ReleaseYear.IsNil() {
		if movie, err := r.tmdbClient.SearchMovie(ctx, tmdb.SearchMovieParams{
			Title:                tc.Title,
			Year:                 tc.ReleaseYear,
			IncludeAdult:         true,
			LevenshteinThreshold: 5,
		}); err == nil {
			if err := tc.SetContent(movie); err != nil {
				return model.TorrentContent{}, err
			}
			return tc, nil
		} else if !errors.Is(err, tmdb.ErrNotFound) {
			return model.TorrentContent{}, err
		}
	}
	return model.TorrentContent{}, classifier.ErrNoMatch
}

func (r videoResolver) resolveTvShow(ctx context.Context, tc model.TorrentContent) (model.TorrentContent, error) {
	externalIds := tc.ExternalIds.OrderedEntries()
	if len(externalIds) > 0 {
		for _, id := range externalIds {
			if tvShow, err := r.tmdbClient.GetTvShowByExternalId(ctx, id.Key, id.Value); err == nil {
				if err := tc.SetContent(tvShow); err != nil {
					return model.TorrentContent{}, err
				}
				return tc, nil
			} else if !errors.Is(err, tmdb.ErrNotFound) {
				return model.TorrentContent{}, err
			}
		}
	} else {
		if tvShow, err := r.tmdbClient.SearchTvShow(ctx, tmdb.SearchTvShowParams{
			Name:                 tc.Title,
			FirstAirDateYear:     tc.ReleaseYear,
			IncludeAdult:         true,
			LevenshteinThreshold: 5,
		}); err == nil {
			if err := tc.SetContent(tvShow); err != nil {
				return model.TorrentContent{}, err
			}
			return tc, nil
		} else if !errors.Is(err, tmdb.ErrNotFound) {
			return model.TorrentContent{}, err
		}
	}
	return model.TorrentContent{}, classifier.ErrNoMatch
}
