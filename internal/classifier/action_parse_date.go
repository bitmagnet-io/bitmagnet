package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/parsers"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const parseDateName = "parse_date"

type parseDateAction struct{}

func (parseDateAction) name() string {
	return parseDateName
}

var parseDateSpec = json_spec.Literal[string]{
	Literal:     parseDateName,
	Description: "Try to parse a date from the name of the current torrent",
}

func (parseDateAction) compile(ctx compilerContext) (action, error) {
	if _, err := parseDateSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			parsed := parsers.ParseDate(ctx.result.Torrent.Name)
			if parsed.IsNil() {
				return ctx.result, classification.ErrUnmatched
			}
			cl := ctx.result
			cl.Date = parsed
			return cl, nil
		},
	}, nil
}

func (parseDateAction) JSONSchema() json_schema.JSONSchema {
	return parseDateSpec.JSONSchema()
}
