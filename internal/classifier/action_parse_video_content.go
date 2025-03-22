package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/parsers"
)

const parseVideoContentName = "parse_video_content"

type parseVideoContentAction struct{}

func (parseVideoContentAction) name() string {
	return parseVideoContentName
}

var parseVideoContentPayloadSpec = payloadLiteral[string]{
	literal:     parseVideoContentName,
	description: "Parse video-related attributes from the name of the current torrent",
}

func (parseVideoContentAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := parseVideoContentPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			parsed, err := parsers.ParseVideoContent(ctx.torrent, ctx.result)
			cl := ctx.result
			if err != nil {
				return cl, err
			}
			cl.Merge(parsed)
			return cl, nil
		},
	}, nil
}

func (parseVideoContentAction) JSONSchema() JSONSchema {
	return parseVideoContentPayloadSpec.JSONSchema()
}
