package httpserver

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
)

const acceptLanguageHeader = "Accept-Language"

func AcceptLanguageFromContext(ctx context.Context) []string {
	if ginContext, ok := httpserver.GinContextFromContext(ctx); ok {
		r := ginContext.Request.Header.Values(acceptLanguageHeader)
		return r
	}

	return nil
}

func NewLocalizerFromContext(ctx context.Context, bundle *i18n.Bundle) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, AcceptLanguageFromContext(ctx)...)
}
