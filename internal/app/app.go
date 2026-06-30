package app

import (
	"context"
	"log/slog"
	"os"

	grpcapp "github.com/barnigator/sso/internal/app/grpc"
	"github.com/barnigator/sso/internal/auth/repository/sqlite"
	slogpretty "github.com/barnigator/sso/pkg/logger/handlers/slogprety"

	"time"

	"github.com/barnigator/sso/internal/auth/usecase"
	"github.com/barnigator/sso/internal/infrastructure/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    *sqlite.Storage
	Log        *slog.Logger
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := usecase.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
		Log:        log,
	}
}

func Run() (*App, error) {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("cfg", cfg))

	application := New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go func() {
		application.GRPCServer.MustRun()
	}()

	return application, nil

}

func (a *App) Stop(ctx context.Context) {
	const fn = "app.Stop"
	a.Log.With(slog.String("fn", fn)).Info("stopping application")

	a.GRPCServer.Stop(ctx)

	if err := a.Storage.Close(); err != nil {
		a.Log.Error("failed to close db", "error", err)
	}

	a.Log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
