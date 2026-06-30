package main

import (
	"context"
	"os/signal"
	"time"

	"github.com/barnigator/sso/internal/app"

	"syscall"
)

func main() {
	application, err := app.Run()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application.Stop(shutDownCtx)
}

// go run cmd/sso/main.go --config=./config/local.yaml
