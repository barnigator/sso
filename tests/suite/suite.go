package suite

import (
	"context"
	"database/sql"
	"net"
	"strconv"
	"testing"

	"github.com/barnigator/sso/internal/infrastructure/config"
	"github.com/stretchr/testify/require"

	ssov1 "github.com/barnigator/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()

	cfg := config.MustLoadByPath("../config/local.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	prepareDB(t, cfg.StoragePath)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}

}

func prepareDB(t *testing.T, dsn string) {
	t.Helper()

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		TRUNCATE TABLE users, apps, admins
		RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO apps(id, name, secret)
		VALUES (1, 'test-app', 'test-secret');
	`)
	require.NoError(t, err)
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
