package jwt

import (
	"errors"

	"github.com/barnigator/sso/internal/auth/domain"
	jwtgo "github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID   int64           `json:"uid"`
	Email    string          `json:"email"`
	Role     domain.UserRole `json:"role"`
	IsActive bool            `json:"is_active"`
	AppID    int             `json:"app_id"`

	jwtgo.RegisteredClaims
}

func ValidateToken(rawToken string, secret []byte) (Claims, error) {
	var claims Claims

	if len(secret) == 0 {
		return Claims{}, ErrInvalidToken
	}

	parser := jwtgo.NewParser(
		jwtgo.WithValidMethods([]string{jwtgo.SigningMethodHS256.Alg()}),
		jwtgo.WithExpirationRequired(),
	)

	token, err := parser.ParseWithClaims(rawToken, &claims, func(_ *jwtgo.Token) (any, error) {
		return secret, nil
	})
	if err != nil || token == nil || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}
