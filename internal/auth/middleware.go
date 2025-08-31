package auth

import (
	"net/http"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserContextKey      = "user"
)

type AuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
	OptionalAuth() gin.HandlerFunc
}

type authMiddleware struct {
	jwtService  JWTService
	authService AuthService
}

func NewAuthMiddleware(jwtService JWTService, authService AuthService) AuthMiddleware {
	return &authMiddleware{
		jwtService:  jwtService,
		authService: authService,
	}
}

func (a *authMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user, err := a.authService.GetUserByID(c, claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}

func (a *authMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		user, err := a.authService.GetUserByID(c, claims.UserID)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserContextKey, user)
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

// GetCurrentUser retrieves the current authenticated user from the Gin context
func GetCurrentUser(c *gin.Context) (*model.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}

	authUser, ok := user.(*model.User)
	if !ok {
		return nil, false
	}

	return authUser, true
}
