package webui

import (
	"errors"
	"io/fs"
	"net/http"

	docs "github.com/bitmagnet-io/bitmagnet/bitmagnet.io"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	FS FS
}

type FS fs.FS

var (
	Ref = http_server.Ref.MustSub("docs")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs an instance of the documentation website at the /docs endpoint"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithFxOption[deps](
			fx.Provide(
				func() (FS, error) {
					return fs.Sub(docs.StaticFS(), "_site")
				},
			),
		),
		builder.WithGinOption(
			Ref,
			0,
			func(deps deps) gin.OptionFunc {
				return func(e *gin.Engine) {
					e.StaticFS("/docs", wrappedFs{http.FS(deps.FS)})
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
