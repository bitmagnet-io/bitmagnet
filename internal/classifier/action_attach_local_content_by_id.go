package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const attachLocalContentByIDName = "attach_local_content_by_id"

type attachLocalContentByIDAction struct{}

func (attachLocalContentByIDAction) name() string {
	return attachLocalContentByIDName
}

var attachLocalContentByIDPayloadSpec = payloadLiteral[string]{
	literal:     attachLocalContentByIDName,
	description: "Use the torrent hint to attach locally stored content by ID",
}

func (attachLocalContentByIDAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentByIDPayloadSpec.Unmarshal(ctx); err != nil {
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

func (attachLocalContentByIDAction) JSONSchema() JSONSchema {
	return attachLocalContentByIDPayloadSpec.JSONSchema()
}
