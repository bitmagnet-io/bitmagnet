package classifier

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/parsers"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
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
				ctx.logger.Info("error")
				return cl, err
			}

			nulldate := model.Date{ Year: 0, Month: 0, Day: 0 }
			var mparsed map[string]any
			jparsed, _ := json.Marshal(parsed)
			json.Unmarshal(jparsed, &mparsed)
			for k, v := range mparsed {
				arr, ok := v.([]any)
				if (ok && len(arr) == 0) || v == nil || v == 0 || v == "" || v == "0001-01-01T00:00:00Z" || v == "0000000000000000000000000000000000000000" || v  == nulldate {
					delete(mparsed, k)
				}
			}

			ctx.logger.Infow("result", "parsed", mparsed)
			cl.Merge(parsed)
			return cl, nil
		},
	}, nil
}

func (parseVideoContentAction) JSONSchema() JSONSchema {
	return parseVideoContentPayloadSpec.JSONSchema()
}
