package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/parsers"
)

const parseDateName = "parse_date"

type parseDateAction struct{}

func (parseDateAction) name() string {
	return parseDateName
}

var parseDatePayloadSpec = payloadLiteral[string]{
	literal:     parseDateName,
	description: "Try to parse a date from the name of the current torrent",
}

func (parseDateAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := parseDatePayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			parsed := parsers.ParseDate(ctx.torrent.Name)
			if parsed.IsNil() {
				return ctx.result, classification.ErrUnmatched
			}
			cl := ctx.result
			cl.Date = parsed
			return cl, nil
		},
	}, nil
}

func (parseDateAction) JSONSchema() JSONSchema {
	return parseDatePayloadSpec.JSONSchema()
}
