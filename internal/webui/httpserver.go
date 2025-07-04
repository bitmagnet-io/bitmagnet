package webui

import (
  "errors"
  "io/fs"
  "net/http"

  "github.com/bitmagnet-io/bitmagnet/internal/httpserver"
  "github.com/bitmagnet-io/bitmagnet/internal/logging"
  "github.com/bitmagnet-io/bitmagnet/webui"
  "github.com/gin-gonic/gin"
)

const Namespace = "webui"

func New(logger logging.Logger) httpserver.Option {
  return &builder{
    logger: logger,
  }
}

type builder struct {
  logger logging.Logger
}

func (*builder) Key() string {
  return Namespace
}

func (b *builder) Apply(e *gin.Engine) {
  webuiFS := webui.StaticFS()

  appRoot, appRootErr := fs.Sub(webuiFS, "dist/bitmagnet/browser")
  if appRootErr != nil {
    b.logger.Errorf(
      "the webui app root directory is missing; run `npm run build` within the `webui` folder: %v",
      appRootErr)
  }

  e.StaticFS("/webui", wrappedFs{http.FS(appRoot)})
  e.GET("/", func(c *gin.Context) {
    c.Redirect(301, "/webui")
  })
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
