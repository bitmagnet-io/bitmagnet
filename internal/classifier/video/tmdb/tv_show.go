package tmdb

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	tmdb "github.com/cyruzin/golang-tmdb"
	"strconv"
)

type TvShowClient interface {
	SearchTvShow(ctx context.Context, p SearchTvShowParams) (model.Content, error)
	GetTvShowByExternalId(ctx context.Context, source, id string) (model.Content, error)
}

type SearchTvShowParams struct {
	Name                 string
	FirstAirDateYear     model.Year
	IncludeAdult         bool
	LevenshteinThreshold uint
}

func (c *client) SearchTvShow(ctx context.Context, p SearchTvShowParams) (tvShow model.Content, err error) {
	if localResult, localErr := c.searchTvShowLocal(ctx, p); localErr == nil {
		return localResult, nil
	} else if !errors.Is(localErr, ErrNotFound) {
		err = localErr
		return
	}
	return c.searchTvShowTmdb(ctx, p)
}

func (c *client) searchTvShowLocal(ctx context.Context, p SearchTvShowParams) (tvShow model.Content, err error) {
	options := []query.Option{
		query.Where(search.ContentTypeCriteria(model.ContentTypeTvShow)),
		query.QueryString(p.Name),
		query.OrderByQueryStringRank(),
		query.Limit(5),
	}
	if !p.FirstAirDateYear.IsNil() {
		options = append(options, query.Where(search.ContentReleaseDateCriteria(model.NewDateRangeFromYear(p.FirstAirDateYear))))
	}
	result, searchErr := c.s.Content(
		ctx,
		options...,
	)
	if searchErr != nil {
		err = searchErr
		return
	}
	for _, item := range result.Items {
		candidates := []string{item.Title}
		if item.OriginalTitle.Valid {
			candidates = append(candidates, item.OriginalTitle.String)
		}
		if levenshteinCheck(p.Name, candidates, p.LevenshteinThreshold) {
			return item.Content, nil
		}
	}
	err = ErrNotFound
	return
}

func (c *client) searchTvShowTmdb(ctx context.Context, p SearchTvShowParams) (tvShow model.Content, err error) {
	urlOptions := make(map[string]string)
	if !p.FirstAirDateYear.IsNil() {
		urlOptions["first_air_date_year"] = strconv.Itoa(int(p.FirstAirDateYear))
	}
	if p.IncludeAdult {
		urlOptions["include_adult"] = "true"
	}
	searchResult, searchErr := c.c.GetSearchTVShow(
		p.Name,
		urlOptions,
	)
	if searchErr != nil {
		err = searchErr
		return
	}
	for _, item := range searchResult.Results {
		if levenshteinCheck(p.Name, []string{item.Name, item.OriginalName}, p.LevenshteinThreshold) {
			return c.GetTvShowByExternalId(ctx, SourceTmdb, strconv.Itoa(int(item.ID)))
		}
	}
	err = ErrNotFound
	return
}

func (c *client) GetTvShowByExternalId(ctx context.Context, source, id string) (tvShow model.Content, err error) {
	searchResult, searchErr := c.s.Content(ctx, query.Where(search.ContentIdentifierCriteria(model.ContentRef{
		Type:   model.ContentTypeTvShow,
		Source: source,
		ID:     id,
	})), query.Limit(1))
	if searchErr != nil {
		return model.Content{}, searchErr
	}
	if len(searchResult.Items) > 0 {
		return searchResult.Items[0].Content, nil
	}
	if source == SourceTmdb {
		intId, idErr := strconv.Atoi(id)
		if idErr != nil {
			err = idErr
			return
		}
		return c.getTvShowByTmdbId(ctx, intId)
	}
	externalSource, externalId, externalSourceErr := getExternalSource(source, id)
	if externalSourceErr != nil {
		err = externalSourceErr
		return
	}
	byIdResult, byIdErr := c.c.GetFindByID(externalId, map[string]string{
		"external_source": externalSource,
	})
	if byIdErr != nil {
		err = byIdErr
		return
	}
	if len(byIdResult.TvResults) == 0 {
		err = ErrNotFound
		return
	}
	return c.getTvShowByTmdbId(ctx, int(byIdResult.TvResults[0].ID))
}

func (c *client) getTvShowByTmdbId(ctx context.Context, id int) (tvShow model.Content, err error) {
	d, getDetailsErr := c.c.GetTVDetails(id, map[string]string{
		"append_to_response": "external_ids",
	})
	if getDetailsErr != nil {
		err = getDetailsErr
		return
	}
	return TvShowDetailsToTvShowModel(*d)
}

func TvShowDetailsToTvShowModel(details tmdb.TVDetails) (movie model.Content, err error) {
	firstAirDate := model.Date{}
	if details.FirstAirDate != "" {
		parsedDate, parseDateErr := model.NewDateFromIsoString(details.FirstAirDate)
		if parseDateErr != nil {
			err = parseDateErr
			return
		}
		firstAirDate = parsedDate
	}
	var collections []model.ContentCollection
	for _, genre := range details.Genres {
		collections = append(collections, model.ContentCollection{
			Type:   "genre",
			Source: SourceTmdb,
			ID:     strconv.Itoa(int(genre.ID)),
			Name:   genre.Name,
		})
	}
	var attributes []model.ContentAttribute
	if details.IMDbID != "" {
		attributes = append(attributes, model.ContentAttribute{
			Source: "imdb",
			Key:    "id",
			Value:  details.IMDbID,
		})
	}
	if details.TVDBID != 0 {
		attributes = append(attributes, model.ContentAttribute{
			Source: "tvdb",
			Key:    "id",
			Value:  strconv.Itoa(int(details.TVDBID)),
		})
	}
	releaseYear := firstAirDate.Year
	if details.PosterPath != "" {
		attributes = append(attributes, model.ContentAttribute{
			Source: "tmdb",
			Key:    "poster_path",
			Value:  details.PosterPath,
		})
	}
	if details.BackdropPath != "" {
		attributes = append(attributes, model.ContentAttribute{
			Source: "tmdb",
			Key:    "backdrop_path",
			Value:  details.BackdropPath,
		})
	}
	return model.Content{
		Type:             model.ContentTypeTvShow,
		Source:           SourceTmdb,
		ID:               strconv.Itoa(int(details.ID)),
		Title:            details.Name,
		ReleaseDate:      firstAirDate,
		ReleaseYear:      releaseYear,
		OriginalLanguage: model.ParseLanguage(details.OriginalLanguage),
		OriginalTitle:    model.NewNullString(details.OriginalName),
		Overview: model.NullString{
			String: details.Overview,
			Valid:  details.Overview != "",
		},
		Popularity:  model.NewNullFloat32(details.Popularity),
		VoteAverage: model.NewNullFloat32(details.VoteAverage),
		VoteCount:   model.NewNullUint(uint(details.VoteCount)),
		Collections: collections,
		Attributes:  attributes,
	}, nil
}
