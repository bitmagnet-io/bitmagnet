package registry

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

type instance struct {
	ref        ref.Ref
	enabled    bool
	dependsOn  []ref.Ref
	requiredBy []ref.Ref
	localizer  plugin.LocalizedContentProvider
}

var _ plugin.Instance = (*instance)(nil)

func (i *instance) Ref() ref.Ref {
	return i.ref
}

func (i *instance) Enabled() bool {
	return i.enabled
}

func (i *instance) DependsOn() []ref.Ref {
	return i.dependsOn
}

func (i *instance) RequiredBy() []ref.Ref {
	return i.requiredBy
}

func (i *instance) LocalizedContent(
	ctx context.Context,
	acceptLanguage ...string,
) (plugin.LocalizedContent, error) {
	return i.localizer.LocalizedContent(ctx, acceptLanguage...)
}
