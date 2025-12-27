package classifier

import (
	"strconv"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const attachTMDBContentByIDName = "attach_tmdb_content_by_id"

type attachTMDBContentByIDAction struct{}

func (attachTMDBContentByIDAction) name() string {
	return attachTMDBContentByIDName
}

var attachTMDBContentByIDPayloadSpec = payloadLiteral[string]{
	literal:     attachTMDBContentByIDName,
	description: "Use the torrent hint to attach content from the TMDB API by ID",
}

func (attachTMDBContentByIDAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachTMDBContentByIDPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			var ref model.ContentRef
			maybeRef := ctx.torrent.Hint.ContentRef()
			if !maybeRef.Valid {
				ctx.logger.Info("hint missing or invalid")
				return cl, classification.ErrUnmatched
			}
			ref = maybeRef.Val
			if cl.ContentType.Valid {
				ref.Type = cl.ContentType.ContentType
			}
			var tmdbID int64
			switch ref.Source {
			case model.SourceTmdb:
				id, err := strconv.Atoi(ref.ID)
				if err != nil {
					ctx.logger.Info("invalid tmdb id")
					return cl, classification.ErrUnmatched
				}
				tmdbID = int64(id)
			default:
				id, err := ctx.tmdbGetTMDBIDByExternalID(ref)
				if err != nil {
					ctx.logger.Info("failed to get tmdb id by external id")
					return cl, err
				}
				tmdbID = id
			}
			var content *model.Content
			switch ref.Type {
			case model.ContentTypeMovie, model.ContentTypeXxx:
				c, err := ctx.tmdbGetMovieByTMDBID(tmdbID)
				if err != nil {
					ctx.logger.Infow("movie not found", "id", tmdbID)
					return cl, err
				}
				ctx.logger.Infow(
					"movie",
					"id", tmdbID,
					"title", c.Title,
					"year", c.ReleaseYear.String())
				content = &c
			case model.ContentTypeTvShow:
				c, err := ctx.tmdbGetTVShowByTMDBID(tmdbID)
				if err != nil {
					ctx.logger.Infow("tv show not found", "id", tmdbID)
					return cl, err
				}
				ctx.logger.Infow(
					"tv show",
					"id", tmdbID,
					"title", c.Title,
					"year", c.ReleaseYear.String())
				content = &c
			default:
				ctx.logger.Info("invalid content type")
				return cl, classification.ErrUnmatched
			}
			cl.AttachContent(content)
			return cl, nil
		},
	}, nil
}

func (attachTMDBContentByIDAction) JSONSchema() JSONSchema {
	return attachTMDBContentByIDPayloadSpec.JSONSchema()
}
