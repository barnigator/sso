package deps

import (
	"context"

	"github.com/barnigator/sso/internal/auth/domain"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	GetUser(ctx context.Context, userID int64) (domain.User, error)
}
