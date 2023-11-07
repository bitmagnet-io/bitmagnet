package video

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"go.uber.org/fx"
	"strconv"
	"strings"
)

type Params struct {
	fx.In
	TmdbClient tmdb.Client
}

type Result struct {
	fx.Out
	Resolver resolver.SubResolver `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: videoResolver{
			config:     resolver.SubResolverConfig{Key: "video", Priority: 1},
			tmdbClient: p.TmdbClient,
		},
	}
}

type videoResolver struct {
	config     resolver.SubResolverConfig
	tmdbClient tmdb.Client
}

func (r videoResolver) Config() resolver.SubResolverConfig {
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
	return model.TorrentContent{}, resolver.ErrNoMatch
}

func (r videoResolver) resolveMovie(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {

	// The movie Onlyfans an horror story always false positive with onlyfans adult site
	titleLower := strings.ToLower(content.Torrent.Name)
	if strings.Contains(titleLower, "onlyfans") && strings.Contains(titleLower, "xxx") {
		output := content
		output.ContentType.Valid = true
		output.ContentType.ContentType = model.ContentTypeXxx
		return output, nil
	}

	externalIds := content.ExternalIds.OrderedEntries()
	if len(externalIds) > 0 {
		for _, id := range externalIds {
			if movie, err := r.tmdbClient.GetMovieByExternalId(ctx, id.Key, id.Value); err == nil {
				content.Content = movie
				return postEnrich(content), nil
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
			return postEnrich(content), nil
		} else if !errors.Is(err, tmdb.ErrNotFound) {
			return model.TorrentContent{}, err
		}
	}
	return model.TorrentContent{}, resolver.ErrNoMatch
}

func (r videoResolver) resolveTvShow(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	externalIds := content.ExternalIds.OrderedEntries()
	if len(externalIds) > 0 {
		for _, id := range externalIds {
			if tvShow, err := r.tmdbClient.GetTvShowByExternalId(ctx, id.Key, id.Value); err == nil {
				content.Content = tvShow
				return postEnrich(content), nil
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
			return postEnrich(content), nil
		} else if !errors.Is(err, tmdb.ErrNotFound) {
			return model.TorrentContent{}, err
		}
	}
	return model.TorrentContent{}, resolver.ErrNoMatch
}

func postEnrich(tc model.TorrentContent) model.TorrentContent {
	c := tc.Content
	contentType := c.Type
	if c.Adult.Valid && c.Adult.Bool {
		contentType = model.ContentTypeXxx
	}
	tc.ContentType = model.NewNullContentType(contentType)
	tc.ContentSource = model.NewNullString(c.Source)
	tc.ContentID = model.NewNullString(c.ID)
	titleParts := []string{c.Title}
	searchStringParts := []string{c.Title}
	if c.OriginalTitle.Valid && c.Title != c.OriginalTitle.String {
		titleParts = append(titleParts, fmt.Sprintf("/ %s", c.OriginalTitle.String))
		searchStringParts = append(searchStringParts, c.OriginalTitle.String)
	}
	if !c.ReleaseDate.IsNil() {
		tc.ReleaseDate = c.ReleaseDate
		tc.ReleaseYear = c.ReleaseDate.Year
	}
	if !tc.ReleaseYear.IsNil() {
		titleParts = append(titleParts, fmt.Sprintf("(%d)", tc.ReleaseYear))
		searchStringParts = append(searchStringParts, strconv.Itoa(int(tc.ReleaseYear)))
	}
	if len(tc.Languages) == 0 && c.OriginalLanguage.Valid {
		tc.Languages = model.Languages{c.OriginalLanguage.Language: struct{}{}}
	}
	if len(tc.Episodes) > 0 {
		titleParts = append(titleParts, tc.Episodes.String())
	}
	searchStringParts = append(searchStringParts, additionalSearchStringParts(tc)...)
	for _, c := range c.Collections {
		if c.Type == "genre" {
			searchStringParts = append(searchStringParts, c.Name)
		}
	}
	tc.Title = strings.Join(titleParts, " ")
	tc.SearchString = strings.Join(searchStringParts, " ")
	return tc
}

func additionalSearchStringParts(content model.TorrentContent) []string {
	var searchStringParts []string
	if content.VideoResolution.Valid {
		searchStringParts = append(searchStringParts, string(content.VideoResolution.VideoResolution))
	}
	if content.VideoSource.Valid {
		searchStringParts = append(searchStringParts, content.VideoSource.VideoSource.String())
	}
	if content.VideoCodec.Valid {
		searchStringParts = append(searchStringParts, string(content.VideoCodec.VideoCodec))
	}
	if content.VideoModifier.Valid {
		searchStringParts = append(searchStringParts, string(content.VideoModifier.VideoModifier))
	}
	if content.ReleaseGroup.Valid {
		searchStringParts = append(searchStringParts, content.ReleaseGroup.String)
	}
	searchStringParts = append(searchStringParts, regex.NormalizeString(content.Torrent.Name))
	return searchStringParts
}
