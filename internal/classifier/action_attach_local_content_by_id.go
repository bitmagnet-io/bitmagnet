package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const attachLocalContentByIdName = "attach_local_content_by_id"

type attachLocalContentByIdAction struct{}

func (attachLocalContentByIdAction) name() string {
	return attachLocalContentByIdName
}

var attachLocalContentByIdPayloadSpec = payloadLiteral[string]{
	literal:     attachLocalContentByIdName,
	description: "Use the torrent hint to attach locally stored content by ID",
}

func (a attachLocalContentByIdAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentByIdPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if ctx.torrent.Hint.IsNil() || !ctx.torrent.Hint.ContentSource.Valid {
				return cl, classification.ErrUnmatched
			}
			content, err := ctx.search.ContentByID(ctx, model.ContentRef{
				Type:   ctx.torrent.Hint.ContentType,
				Source: ctx.torrent.Hint.ContentSource.String,
				ID:     ctx.torrent.Hint.ContentID.String,
			})
			if err != nil {
				return cl, err
			}
			cl.AttachContent(&content)
			return cl, nil
		},
	}, nil
}

func (a attachLocalContentByIdAction) JSONSchema() JSONSchema {
	return attachLocalContentByIdPayloadSpec.JSONSchema()
}
