package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/processor/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/workflow/parsers"
)

const parseDateName = "parse_date"

type parseDateAction struct{}

func (parseDateAction) Name() string {
	return parseDateName
}

var parseDatePayloadSpec = payloadLiteral[string]{parseDateName}

func (parseDateAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := parseDatePayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			parsed := parsers.ParseDate(ctx.torrent.Name)
			if parsed.IsNil() {
				return ctx.result, classification.ErrNoMatch
			}
			cl := ctx.result
			cl.Date = parsed
			return cl, nil
		},
	}, nil
}
