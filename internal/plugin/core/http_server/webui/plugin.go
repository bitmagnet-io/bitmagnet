package webui

import (
	"errors"
	"io/fs"
	"net/http"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/webui"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type (
	config struct{}
	deps   struct {
		fx.In
		FS FS
	}
)

type FS fs.FS

var (
	Ref = http_server.Ref.MustSub("webui")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithFxOption[config, deps](
			fx.Provide(
				func() (FS, error) {
					return fs.Sub(webui.StaticFs, "dist/bitmagnet/browser")
				},
			),
		),
		builder.WithGinOption(
			Ref,
			func(_ config, deps deps) gin.OptionFunc {
				return func(e *gin.Engine) {
					e.StaticFS("/webui", wrappedFs{http.FS(deps.FS)})
					e.GET("/", func(c *gin.Context) {
						c.Redirect(301, "/webui")
					})
				}
			},
		),
	)
)

type wrappedFs struct {
	http.FileSystem
}

func (w wrappedFs) Open(name string) (http.File, error) {
	f, err := w.FileSystem.Open(name)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return w.FileSystem.Open("/index.html")
	}

	return f, err
}
