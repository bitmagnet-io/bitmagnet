package http_auth

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/identity"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	IdentityKey         = "identity"
)

type Middleware interface {
	AttachAuth() gin.HandlerFunc
}

type authMiddleware struct {
	authenticator identity.Authenticator
}

func NewMiddleware(authenticator identity.Authenticator) Middleware {
	return &authMiddleware{
		authenticator: authenticator,
	}
}

func (a *authMiddleware) AttachAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, matched, err := a.authenticator.Authenticate(c, extractToken(c))

		if err == nil && matched {
			c.Set(IdentityKey, identity)
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, BearerPrefix)
}

// GetIdentity retrieves the current authenticated identity from the Gin context
func GetIdentity(c *gin.Context) (identity.Identity, bool) {
	raw, exists := c.Get(IdentityKey)
	if !exists {
		return nil, false
	}

	identity, ok := raw.(identity.Identity)
	if !ok {
		return nil, false
	}

	return identity, true
}
