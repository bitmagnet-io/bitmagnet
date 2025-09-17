package jwt_test

import (
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	jwtService := jwt.NewService("test-secret-key", jwt.Duration(time.Minute))

	userID := 123
	username := "testuser"

	// Generate token
	token, err := jwtService.GenerateToken(userID, username)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Validate token
	claims, err := jwtService.ValidateToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, "bitmagnet", claims.Issuer)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestJWTService_InvalidToken(t *testing.T) {
	jwtService := jwt.NewService("test-secret-key", jwt.Duration(time.Minute))

	// Test with invalid token
	_, err := jwtService.ValidateToken("invalid-token")
	assert.Error(t, err)

	// Test with empty token
	_, err = jwtService.ValidateToken("")
	assert.Error(t, err)
}

func TestJWTService_TokenWithDifferentSecret(t *testing.T) {
	jwtService1 := jwt.NewService("secret1", jwt.Duration(time.Minute))
	jwtService2 := jwt.NewService("secret2", jwt.Duration(time.Minute))

	userID := 123
	username := "testuser"

	// Generate token with first service
	token, err := jwtService1.GenerateToken(userID, username)
	require.NoError(t, err)

	// Try to validate with second service (different secret)
	_, err = jwtService2.ValidateToken(token)
	assert.Error(t, err)
}
