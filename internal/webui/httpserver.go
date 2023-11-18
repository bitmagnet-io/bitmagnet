package webui

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/webui"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
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
	appRoot, appRootErr := fs.Sub(webuiFS, "dist/bitmagnet")
	if appRootErr != nil {
		b.logger.Errorf("the webui app root directory is missing; run `ng build` within the `webui` folder: %v", appRootErr)
		return nil
	}
	appRootFS := http.FS(appRoot)
	if walkErr := fs.WalkDir(appRoot, ".", func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() {
			e.StaticFileFS(fmt.Sprintf("/%s", path), fmt.Sprintf("./%s", path), appRootFS)
		}
		return nil
	}); walkErr != nil {
		return fmt.Errorf("failed to walk app root: %w", walkErr)
	}
	// serving index.html from "/" using e.StaticFileFS seems to result in a redirect loop:
	index, indexErr := appRootFS.Open("index.html")
	if indexErr != nil {
		return fmt.Errorf("failed to open index.html: %w", indexErr)
	}
	indexBytes, indexBytesErr := io.ReadAll(index)
	if indexBytesErr != nil {
		return fmt.Errorf("failed to read index.html: %w", indexBytesErr)
	}
	e.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		_, _ = c.Writer.Write(indexBytes)
	})
	return nil
}
