package workflow

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/classification"
)

const attachLocalContentBySearchName = "attach_local_content_by_search"

type attachLocalContentBySearchAction struct {
	searchAction
}

func (attachLocalContentBySearchAction) Name() string {
	return attachLocalContentBySearchName
}

var attachLocalContentBySearchPayloadSpec = payloadLiteral[string]{attachLocalContentBySearchName}

func (a attachLocalContentBySearchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentBySearchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if !cl.ContentType.Valid || !cl.BaseTitle.Valid {
				return cl, classification.ErrNoMatch
			}
			options := []query.Option{
				query.Where(search.ContentTypeCriteria(cl.ContentType.ContentType)),
				query.QueryString(fmt.Sprintf("\"%s\"", cl.BaseTitle.String)),
				query.OrderByQueryStringRank(),
				query.Limit(5),
				search.ContentDefaultPreload(),
				search.ContentDefaultHydrate(),
			}
			if !cl.Date.Year.IsNil() {
				options = append(options, query.Where(search.ContentReleaseDateCriteria(model.NewDateRangeFromYear(cl.Date.Year))))
			}
			result, searchErr := a.search.Content(
				ctx,
				options...,
			)
			if searchErr != nil {
				return cl, searchErr
			}
			var content *model.Content
			for _, item := range result.Items {
				candidates := []string{item.Title}
				if item.OriginalTitle.Valid {
					candidates = append(candidates, item.OriginalTitle.String)
				}
				if levenshteinCheck(cl.BaseTitle.String, candidates, levenshteinThreshold) {
					c := item.Content
					content = &c
					break
				}
			}
			if content == nil {
				return cl, classification.ErrNoMatch
			}
			cl.AttachContent(content)
			return cl, nil
		},
	}, nil
}
