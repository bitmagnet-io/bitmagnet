package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
)

const attachLocalContentBySearchName = "attach_local_content_by_search"

type attachLocalContentBySearchAction struct{}

func (attachLocalContentBySearchAction) name() string {
	return attachLocalContentBySearchName
}

var attachLocalContentBySearchPayloadSpec = payloadLiteral[string]{
	literal:     attachLocalContentBySearchName,
	description: "Attempt to attach local content with a search on the torrent name",
}

func (attachLocalContentBySearchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentBySearchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if !cl.ContentType.Valid || !cl.BaseTitle.Valid {
				ctx.logger.Info("invalid content type or base title")
				return cl, classification.ErrUnmatched
			}
			content, err := ctx.search.ContentBySearch(
				ctx.Context,
				cl.ContentType.ContentType,
				cl.BaseTitle.String,
				cl.Date.Year,
			)
			if err != nil {
				ctx.logger.Infow(
					"local search failed",
					"type",cl.ContentType.ContentType.String(),
					"base_title",cl.BaseTitle.String,
					"date",cl.Date.IsoDateString())
				return cl, err
			}
			ctx.logger.Infow(
				"local search succeeded",
				"base_title",cl.BaseTitle.String,
				"date",cl.Date.IsoDateString(),
				"id",content.ID,
				"type",content.Type.String(),
				"title",content.Title,
				"year",content.ReleaseYear.String())
			cl.AttachContent(&content)
			return cl, nil
		},
	}, nil
}

func (attachLocalContentBySearchAction) JSONSchema() JSONSchema {
	return attachLocalContentBySearchPayloadSpec.JSONSchema()
}
