package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const attachLocalContentBySearchName = "attach_local_content_by_search"

type attachLocalContentBySearchAction struct{}

func (attachLocalContentBySearchAction) name() string {
	return attachLocalContentBySearchName
}

var attachLocalContentBySearchSpec = json_spec.Literal[string]{
	Literal:     attachLocalContentBySearchName,
	Description: "Attempt to attach local content with a search on the torrent name",
}

func (attachLocalContentBySearchAction) compile(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentBySearchSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if !cl.ContentType.Valid || !cl.BaseTitle.Valid {
				return cl, classification.ErrUnmatched
			}
			content, err := ctx.search.ContentBySearch(
				ctx.Context,
				cl.ContentType.ContentType,
				cl.BaseTitle.String,
				cl.Date.Year,
			)
			if err != nil {
				return cl, err
			}
			cl.AttachContent(&content)
			return cl, nil
		},
	}, nil
}

func (attachLocalContentBySearchAction) JSONSchema() json_schema.JSONSchema {
	return attachLocalContentBySearchSpec.JSONSchema()
}
