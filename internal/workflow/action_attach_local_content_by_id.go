package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const attachLocalContentByIdName = "attach_local_content_by_id"

type attachLocalContentByIdAction struct {
	search search.Search
}

func (attachLocalContentByIdAction) Name() string {
	return attachLocalContentByIdName
}

var attachLocalContentByIdPayloadSpec = payloadLiteral[string]{attachLocalContentByIdName}

func (a attachLocalContentByIdAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentByIdPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classifier.Classification, error) {
			cl := ctx.result
			if ctx.torrent.Hint.IsNil() || !ctx.torrent.Hint.ContentSource.Valid {
				return cl, ErrNoMatch
			}
			options := []query.Option{
				query.Where(
					search.ContentTypeCriteria(ctx.torrent.Hint.ContentType),
				),
				search.ContentDefaultPreload(),
				search.ContentDefaultHydrate(),
				query.Limit(1),
			}
			source := ctx.torrent.Hint.ContentSource.String
			id := ctx.torrent.Hint.ContentID.String
			if source == "tmdb" {
				canonicalResult, canonicalErr := a.search.Content(ctx,
					append(options, query.Where(
						search.ContentCanonicalIdentifierCriteria(model.ContentRef{
							Source: source,
							ID:     id,
						}),
					))...,
				)
				if canonicalErr != nil {
					return cl, canonicalErr
				}
				if len(canonicalResult.Items) == 0 {
					return cl, ErrNoMatch
				}
				content := canonicalResult.Items[0].Content
				cl.Content = &content
				return cl, nil
			} else {
				alternativeResult, alternativeErr := a.search.Content(ctx,
					append(options, query.Where(
						search.ContentAlternativeIdentifierCriteria(model.ContentRef{
							Source: source,
							ID:     id,
						}),
					))...,
				)
				if alternativeErr != nil {
					return cl, ErrNoMatch
				}
				if len(alternativeResult.Items) == 0 {
					return cl, ErrNoMatch
				}
				content := alternativeResult.Items[0].Content
				cl.Content = &content
				return cl, nil
			}
		},
	}, nil
}
