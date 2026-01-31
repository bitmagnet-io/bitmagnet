package httpserver

import (
	"context"
	"errors"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bitmagnet-io/bitmagnet/internal/error_registry"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver/batch"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func New(
	schema graphql.ExecutableSchema,
	errorRegistry error_registry.Registry,
	i18n *i18n.Bundle,
) gin.OptionFunc {
	return func(e *gin.Engine) {
		gql := newServer(schema, errorRegistry, i18n)

		e.POST("/graphql", func(c *gin.Context) {
			gql.ServeHTTP(c.Writer, c.Request)
		})

		pg := playground.Handler("GraphQL playground", "/graphql")

		e.GET("/graphql", func(c *gin.Context) {
			pg.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func newServer(
	schema graphql.ExecutableSchema,
	errorRegistry error_registry.Registry,
	i18nBundle *i18n.Bundle,
) *handler.Server {
	srv := handler.New(schema)

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(batch.Transport{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	srv.SetErrorPresenter(newErrorPresenter(errorRegistry, i18nBundle))

	return srv
}

func newErrorPresenter(
	errorRegistry error_registry.Registry,
	i18nBundle *i18n.Bundle,
) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)

		if gqlErr.Extensions == nil {
			gqlErr.Extensions = make(map[string]any)
		}

		if entry, ok := errorRegistry.Identify(err); ok {
			message := entry.Value.Error()

			localizer := NewLocalizerFromContext(ctx, i18nBundle)

			if localized, _ := localizer.LocalizeMessage(&i18n.Message{
				ID: entry.String(),
			}); localized != "" {
				message = localized
			}

			gqlErr.Message = message

			gqlErr.Extensions["code"] = entry.String()
		}

		var hasExtensions HasExtensions
		if errors.As(err, &hasExtensions) {
			for k, v := range hasExtensions.GraphQLExtensions() {
				gqlErr.Extensions[k] = v
			}
		}

		return gqlErr
	}
}

type HasExtensions interface{ GraphQLExtensions() map[string]any }
