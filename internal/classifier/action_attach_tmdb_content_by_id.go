package classifier

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"strconv"
)

const attachTmdbContentByIdName = "attach_tmdb_content_by_id"

type attachTmdbContentByIdAction struct {
	tmdbAction
}

func (attachTmdbContentByIdAction) Name() string {
	return attachTmdbContentByIdName
}

var attachTmdbContentByIdPayloadSpec = payloadLiteral[string]{attachTmdbContentByIdName}

func (a attachTmdbContentByIdAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachTmdbContentByIdPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			var ref model.ContentRef
			if maybeRef := ctx.torrent.Hint.ContentRef(); !maybeRef.Valid {
				return cl, classification.ErrNoMatch
			} else {
				ref = maybeRef.Val
			}
			if cl.ContentType.Valid {
				ref.Type = cl.ContentType.ContentType
			}
			var tmdbId int64
			switch ref.Source {
			case SourceTmdb:
				if id, err := strconv.Atoi(ref.ID); err != nil {
					return cl, classification.ErrNoMatch
				} else {
					tmdbId = int64(id)
				}
			default:
				if id, err := a.getTmdbIdByExternalId(ctx, ref); err != nil {
					return cl, err
				} else {
					tmdbId = id
				}
			}
			var content *model.Content
			switch ref.Type {
			case model.ContentTypeMovie, model.ContentTypeXxx:
				if c, err := a.getMovieByTmbdId(ctx, tmdbId); err != nil {
					return cl, err
				} else {
					content = &c
				}
			case model.ContentTypeTvShow:
				if c, err := a.getTvShowByTmbdId(ctx, tmdbId); err != nil {
					return cl, err
				} else {
					content = &c
				}
			default:
				return cl, classification.ErrNoMatch
			}
			cl.AttachContent(content)
			return cl, nil
		},
	}, nil
}

func (a attachTmdbContentByIdAction) getTmdbIdByExternalId(ctx context.Context, ref model.ContentRef) (int64, error) {
	externalSource, externalId, externalSourceErr := getTmdbExternalSource(ref)
	if externalSourceErr != nil {
		return 0, externalSourceErr
	}
	byIdResult, byIdErr := a.client.FindByID(ctx, tmdb.FindByIDRequest{
		ExternalSource: externalSource,
		ExternalID:     externalId,
	})
	if byIdErr != nil {
		return 0, byIdErr
	}
	switch ref.Type {
	case model.ContentTypeMovie, model.ContentTypeXxx:
		if len(byIdResult.MovieResults) == 0 {
			return 0, classification.ErrNoMatch
		}
		return byIdResult.MovieResults[0].ID, nil
	case model.ContentTypeTvShow:
		if len(byIdResult.TvResults) == 0 {
			return 0, classification.ErrNoMatch
		}
		return byIdResult.TvResults[0].ID, nil
	default:
		return 0, classification.ErrNoMatch
	}
}

func getTmdbExternalSource(ref model.ContentRef) (externalSource string, externalId string, err error) {
	switch {
	case (ref.Type == model.ContentTypeMovie ||
		ref.Type == model.ContentTypeTvShow ||
		ref.Type == model.ContentTypeXxx) &&
		ref.Source == SourceImdb:
		externalSource = "imdb_id"
		externalId = ref.ID
	case ref.Type == model.ContentTypeTvShow && ref.Source == SourceTvdb:
		externalSource = "tvdb_id"
		externalId = ref.ID
	default:
		err = classification.ErrNoMatch
	}
	return
}
