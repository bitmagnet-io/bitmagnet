package webui

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/webui"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

type Params struct {
	fx.In
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		Option: &builder{
			logger: p.Logger.Named("webui"),
		},
	}
}

type builder struct {
	logger *zap.SugaredLogger
}

func (b *builder) Key() string {
	return "webui"
}

func (b *builder) Apply(e *gin.Engine) error {
	webuiFS := webui.StaticFS()
	appRoot, appRootErr := fs.Sub(webuiFS, "dist/bitmagnet/browser")
	if appRootErr != nil {
		b.logger.Errorf("the webui app root directory is missing; run `npm run build` within the `webui` folder: %v", appRootErr)
		return nil
	}
	e.StaticFS("/webui", wrappedFs{http.FS(appRoot)})
	e.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/webui")
	})
	return nil
}

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
