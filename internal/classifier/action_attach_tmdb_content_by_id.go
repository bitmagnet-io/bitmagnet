package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"strconv"
)

const attachTmdbContentByIdName = "attach_tmdb_content_by_id"

type attachTmdbContentByIdAction struct{}

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
			case model.SourceTmdb:
				if id, err := strconv.Atoi(ref.ID); err != nil {
					return cl, classification.ErrNoMatch
				} else {
					tmdbId = int64(id)
				}
			default:
				if id, err := ctx.tmdb_getTmdbIdByExternalId(ref); err != nil {
					return cl, err
				} else {
					tmdbId = id
				}
			}
			var content *model.Content
			switch ref.Type {
			case model.ContentTypeMovie, model.ContentTypeXxx:
				if c, err := ctx.tmdb_getMovieByTmbdId(tmdbId); err != nil {
					return cl, err
				} else {
					content = &c
				}
			case model.ContentTypeTvShow:
				if c, err := ctx.tmdb_getTvShowByTmbdId(tmdbId); err != nil {
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
