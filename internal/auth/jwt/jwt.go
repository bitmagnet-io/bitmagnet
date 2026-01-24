package jwt

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenInvalidClaims = jwt.ErrTokenInvalidClaims

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Service interface {
	Generate(userID int, username string) (string, error)
	Parse(token string) (*Claims, error)
}

type service struct {
	parser        *jwt.Parser
	secretKey     []byte
	tokenDuration time.Duration
}

func NewService(secretKey Secret, duration Duration) Service {
	if secretKey == "" {
		secretKey = Secret(auth.GenerateRandomString(32))
	}

	return &service{
		parser:        jwt.NewParser(),
		secretKey:     []byte(secretKey),
		tokenDuration: time.Duration(duration),
	}
}

func (j *service) Generate(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bitmagnet",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secretKey)
}

func (j *service) Parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(*jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}
