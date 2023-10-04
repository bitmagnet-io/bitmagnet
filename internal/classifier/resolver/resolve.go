package resolver

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"sort"
)

func (r resolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	for _, subResolver := range r.sortedSubResolvers() {
		preEnrichedContent, preEnrichedErr := subResolver.PreEnrich(content)
		if preEnrichedErr != nil {
			return model.TorrentContent{}, preEnrichedErr
		}
		resolvedContent, resolveErr := subResolver.Resolve(ctx, preEnrichedContent)
		if resolveErr == nil {
			return resolvedContent, nil
		}
		if !errors.Is(resolveErr, ErrNoMatch) {
			r.logger.Errorw("error resolving content", "resolver", subResolver.Config().Key, "content", preEnrichedContent, "error", resolveErr)
			return model.TorrentContent{}, resolveErr
		}
	}
	return model.TorrentContent{}, ErrNoMatch
}

func (r resolver) sortedSubResolvers() []SubResolver {
	subResolvers := r.subResolvers
	sort.Slice(subResolvers, func(i, j int) bool {
		return subResolvers[i].Config().Priority < subResolvers[j].Config().Priority
	})
	return subResolvers
}
