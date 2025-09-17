package http_auth

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/jwt"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserContextKey      = "user"
)

type Middleware interface {
	// RequireAuth() gin.HandlerFunc
	AttachAuth() gin.HandlerFunc
}

type authMiddleware struct {
	jwtService    jwt.Service
	apiKeyService api_key.Service
	userService   user.Service
}

func NewMiddleware(jwtService jwt.Service, userService user.Service) Middleware {
	return &authMiddleware{
		jwtService:  jwtService,
		userService: userService,
	}
}

// func (a *authMiddleware) RequireAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := extractToken(c)
// 		if token == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
// 			c.Abort()
// 			return
// 		}

// 		claims, err := a.jwtService.ValidateToken(token)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		user, err := a.userService.Get(c, claims.UserID)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set(UserContextKey, user)
// 		c.Next()
// 	}
// }

func (a *authMiddleware) AttachAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := a.jwtService.ValidateToken(token)
		if err == nil {
			user, err := a.userService.Get(c, claims.UserID)
			if err != nil {
				c.Next()
				return
			}

			c.Set(UserContextKey, user)
			c.Next()
			return
		}
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

// GetCurrentUser retrieves the current authenticated user from the Gin context
func GetCurrentUser(c *gin.Context) (model.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return model.User{}, false
	}

	authUser, ok := user.(model.User)
	if !ok {
		return model.User{}, false
	}

	return authUser, true
}
