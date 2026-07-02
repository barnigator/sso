package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/barnigator/sso/internal/app"
	"github.com/barnigator/sso/internal/infrastructure/config"
)

var testApp *app.App
var testCfg *config.Config

func TestMain(m *testing.M) {
	testCfg = config.MustLoadByPath("../config/test.yaml")

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	var err error
	testApp, err = app.New(
		log,
		testCfg.GRPC.Port,
		testCfg.StoragePath,
		testCfg.TokenTTL,
	)
	if err != nil {
		os.Exit(1)
	}

	go testApp.GRPCServer.MustRun()

	code := m.Run()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testApp.Stop(shutdownCtx)

	os.Exit(code)
}
