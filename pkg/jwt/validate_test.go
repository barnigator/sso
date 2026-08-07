package jwt

import (
	"testing"
	"time"

	"github.com/barnigator/sso/internal/auth/domain"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateToken(t *testing.T) {
	user := domain.User{
		ID:       42,
		Email:    "user@example.com",
		Role:     domain.UserRoleUser,
		IsActive: true,
	}
	app := domain.App{
		ID:     1,
		Secret: "secret",
	}

	token, err := NewToken(user, app, time.Hour)
	require.NoError(t, err)

	claims, err := ValidateToken(token, []byte(app.Secret))
	require.NoError(t, err)

	assert.Equal(t, claims.UserID, user.ID)
	assert.Equal(t, claims.Email, user.Email)
	assert.Equal(t, claims.Role, user.Role)
	assert.True(t, claims.IsActive)
	assert.Equal(t, app.ID, claims.AppID)
	require.NotNil(t, claims.ExpiresAt)
}

func TestValidateToken_Errors(t *testing.T) {

	user := domain.User{
		ID:       42,
		Email:    "user@example.com",
		Role:     domain.UserRoleUser,
		IsActive: true,
	}
	app := domain.App{
		ID:     1,
		Secret: "secret",
	}

	validClaims := Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		AppID:    app.ID,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	expiredClaims := validClaims
	expiredClaims.ExpiresAt = jwtgo.NewNumericDate(time.Now().Add(-time.Minute))

	withoutExpClaims := validClaims
	withoutExpClaims.ExpiresAt = nil

	token, err := NewToken(user, app, time.Hour)
	require.NoError(t, err)

	expiredToken := signTestToken(
		t,
		jwtgo.SigningMethodHS256,
		expiredClaims,
		app.Secret,
	)

	withoutExpToken := signTestToken(
		t,
		jwtgo.SigningMethodHS256,
		withoutExpClaims,
		app.Secret,
	)

	hs512Token := signTestToken(
		t,
		jwtgo.SigningMethodHS512,
		validClaims,
		app.Secret,
	)

	tests := []struct {
		name             string
		token            string
		validationSecret string
	}{
		{
			name:             "wrong secret",
			token:            token,
			validationSecret: "wrong secret",
		},
		{
			name:             "expired token",
			token:            expiredToken,
			validationSecret: app.Secret,
		},
		{
			name:             "token without expiration",
			token:            withoutExpToken,
			validationSecret: app.Secret,
		},
		{
			name:             "wrong signing algorithm",
			token:            hs512Token,
			validationSecret: app.Secret,
		},
		{
			name:             "empty secret",
			token:            token,
			validationSecret: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err = ValidateToken(test.token, []byte(test.validationSecret))
			require.ErrorIs(t, err, ErrInvalidToken)
		})
	}
}

func signTestToken(
	t *testing.T,
	method jwtgo.SigningMethod,
	claims Claims,
	secret string,
) string {
	t.Helper()

	token := jwtgo.NewWithClaims(method, claims)
	signed, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	return signed
}
