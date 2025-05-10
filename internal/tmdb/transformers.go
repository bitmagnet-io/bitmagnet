package tmdb

import (
	"strconv"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func MovieDetailsToMovieModel(details MovieDetailsResponse) (movie model.Content, err error) {
	releaseDate := model.Date{}

	if details.ReleaseDate != "" {
		parsedDate, parseDateErr := model.NewDateFromIsoString(details.ReleaseDate)
		if parseDateErr != nil {
			err = parseDateErr
			return
		}

		releaseDate = parsedDate
	}

	//nolint:prealloc
	var collections []model.ContentCollection

	if details.BelongsToCollection.ID != 0 {
		collections = append(collections, model.ContentCollection{
			Type:   "franchise",
			Source: model.SourceTmdb,
			ID:     strconv.Itoa(int(details.BelongsToCollection.ID)),
			Name:   details.BelongsToCollection.Name,
		})
	}

	for _, genre := range details.Genres {
		collections = append(collections, model.ContentCollection{
			Type:   "genre",
			Source: model.SourceTmdb,
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

	releaseYear := releaseDate.Year

	contentType := model.ContentTypeMovie

	if details.Adult {
		contentType = model.ContentTypeXxx
	}

	return model.Content{
		Type:             contentType,
		Source:           model.SourceTmdb,
		ID:               strconv.Itoa(int(details.ID)),
		Title:            details.Title,
		ReleaseDate:      releaseDate,
		ReleaseYear:      releaseYear,
		Adult:            model.NewNullBool(details.Adult),
		OriginalLanguage: model.ParseLanguage(details.OriginalLanguage),
		OriginalTitle:    model.NewNullString(details.OriginalTitle),
		Overview: model.NullString{
			String: details.Overview,
			Valid:  details.Overview != "",
		},
		Runtime: model.NullUint16{
			Uint16: uint16(details.Runtime),
			Valid:  details.Runtime > 0,
		},
		Popularity:  model.NewNullFloat32(details.Popularity),
		VoteAverage: model.NewNullFloat32(details.VoteAverage),
		VoteCount:   model.NewNullUint(uint(details.VoteCount)),
		Collections: collections,
		Attributes:  attributes,
	}, nil
}

func TvShowDetailsToTvShowModel(details TvDetailsResponse) (movie model.Content, err error) {
	firstAirDate := model.Date{}

	if details.FirstAirDate != "" {
		parsedDate, parseDateErr := model.NewDateFromIsoString(details.FirstAirDate)
		if parseDateErr != nil {
			err = parseDateErr
			return
		}

		firstAirDate = parsedDate
	}

	collections := slice.Map(details.Genres, func(genre Genre) model.ContentCollection {
		return model.ContentCollection{
			Type:   "genre",
			Source: model.SourceTmdb,
			ID:     strconv.Itoa(int(genre.ID)),
			Name:   genre.Name,
		}
	})

	var attributes []model.ContentAttribute

	if details.ExternalIDs.IMDbID != "" {
		attributes = append(attributes, model.ContentAttribute{
			Source: model.SourceImdb,
			Key:    "id",
			Value:  details.ExternalIDs.IMDbID,
		})
	}

	if details.ExternalIDs.TVDBID != 0 {
		attributes = append(attributes, model.ContentAttribute{
			Source: model.SourceTvdb,
			Key:    "id",
			Value:  strconv.Itoa(int(details.ExternalIDs.TVDBID)),
		})
	}

	releaseYear := firstAirDate.Year

	if details.PosterPath != "" {
		attributes = append(attributes, model.ContentAttribute{
			Source: model.SourceTmdb,
			Key:    "poster_path",
			Value:  details.PosterPath,
		})
	}

	if details.BackdropPath != "" {
		attributes = append(attributes, model.ContentAttribute{
			Source: model.SourceTmdb,
			Key:    "backdrop_path",
			Value:  details.BackdropPath,
		})
	}

	return model.Content{
		Type:             model.ContentTypeTvShow,
		Source:           model.SourceTmdb,
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

func ExternalSource(ref model.ContentRef) (externalSource string, externalID string, err error) {
	switch {
	case (ref.Type == model.ContentTypeMovie ||
		ref.Type == model.ContentTypeTvShow ||
		ref.Type == model.ContentTypeXxx) &&
		ref.Source == model.SourceImdb:
		externalSource = "imdb_id"
		externalID = ref.ID
	case ref.Type == model.ContentTypeTvShow && ref.Source == model.SourceTvdb:
		externalSource = "tvdb_id"
		externalID = ref.ID
	default:
		err = classification.ErrUnmatched
	}

	return
}
