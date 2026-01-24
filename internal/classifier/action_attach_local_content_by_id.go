package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const attachLocalContentByIDName = "attach_local_content_by_id"

type attachLocalContentByIDAction struct{}

func (attachLocalContentByIDAction) name() string {
	return attachLocalContentByIDName
}

var attachLocalContentByIDSpec = json_spec.Literal[string]{
	Literal:     attachLocalContentByIDName,
	Description: "Use the torrent hint to attach locally stored content by ID",
}

func (attachLocalContentByIDAction) compile(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentByIDSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if cl.Torrent.Hint.IsNil() || !cl.Torrent.Hint.ContentSource.Valid {
				return cl, classification.ErrUnmatched
			}

			content, err := ctx.search.ContentByID(ctx, model.ContentRef{
				Type:   cl.Torrent.Hint.ContentType,
				Source: cl.Torrent.Hint.ContentSource.String,
				ID:     cl.Torrent.Hint.ContentID.String,
			})
			if err != nil {
				return cl, err
			}

			cl.AttachContent(&content)

			return cl, nil
		},
	}, nil
}

func (attachLocalContentByIDAction) JSONSchema() json_schema.JSONSchema {
	return attachLocalContentByIDSpec.JSONSchema()
}
