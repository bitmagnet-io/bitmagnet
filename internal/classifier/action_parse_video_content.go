package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/parsers"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const parseVideoContentName = "parse_video_content"

type parseVideoContentAction struct{}

func (parseVideoContentAction) name() string {
	return parseVideoContentName
}

var parseVideoContentSpec = json_spec.Literal[string]{
	Literal:     parseVideoContentName,
	Description: "Parse video-related attributes from the name of the current torrent",
}

func (parseVideoContentAction) compile(ctx compilerContext) (action, error) {
	if _, err := parseVideoContentSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			parsed, err := parsers.ParseVideoContent(ctx.result)
			cl := ctx.result
			if err != nil {
				return cl, err
			}
			cl.Merge(parsed)
			return cl, nil
		},
	}, nil
}

func (parseVideoContentAction) JSONSchema() json_schema.JSONSchema {
	return parseVideoContentSpec.JSONSchema()
}
