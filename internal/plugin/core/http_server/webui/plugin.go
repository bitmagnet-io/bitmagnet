package webui

import (
	"errors"
	"io/fs"
	"net/http"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/webui"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	FS              FS
	AnonymousAccess *atomic.Value[rbac.AnonymousAccess]
}

type FS fs.FS

var (
	Ref = http_server.Ref.MustSub("webui")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs the web user interface at the /webui endpoint"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithFxOption[deps](
			fx.Provide(
				func() (FS, error) {
					return fs.Sub(webui.StaticFs, "dist/bitmagnet/browser")
				},
			),
		),
		builder.WithGinOption(
			Ref,
			0,
			func(deps deps) gin.OptionFunc {
				return func(e *gin.Engine) {
					e.StaticFS("/webui", wrappedFs{http.FS(deps.FS)})
					e.GET("/", func(c *gin.Context) {
						c.Redirect(301, "/webui")
					})
				}
			},
		),
		builder.WithAuthObjectActions(
			func(deps) []rbac.ObjectAction {
				return []rbac.ObjectAction{
					authObjectActionPageView,
				}
			},
		),
		builder.WithPermissionProvider(
			func(deps deps) []rbac.Permission {
				permissions := []rbac.Permission{
					rbac.NewPermission(
						rbac.SubjectRole{
							Role: rbac.RoleUser,
						},
						authObjectActionPageView,
					),
				}

				if deps.AnonymousAccess.Get() {
					permissions = append(permissions,
						rbac.NewPermission(
							rbac.SubjectRole{
								Role: rbac.RoleAnon,
							},
							authObjectActionPageView,
						))
				}

				return permissions
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
