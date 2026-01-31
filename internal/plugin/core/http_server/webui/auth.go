package webui

import "github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"

const (
	authNamespace = "webui"
)

var authObjectActionPageView = rbac.NewObjectAction(
	authNamespace,
	"page",
	"view",
)
