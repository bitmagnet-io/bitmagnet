package httpserver

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Schema lazy.Lazy[graphql.ExecutableSchema]
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		Option: &builder{
			schema: p.Schema,
		},
	}
}

type builder struct {
	schema lazy.Lazy[graphql.ExecutableSchema]
}

func (builder) Key() string {
	return "graphql"
}

func (b builder) Apply(e *gin.Engine) error {
	schema, err := b.schema.Get()
	if err != nil {
		return err
	}
	gql := handler.NewDefaultServer(schema)
	e.POST("/graphql", func(c *gin.Context) {
		gql.ServeHTTP(c.Writer, c.Request)
	})
	pg := playground.Handler("GraphQL playground", "/graphql")
	e.GET("/graphql", func(c *gin.Context) {
		pg.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
