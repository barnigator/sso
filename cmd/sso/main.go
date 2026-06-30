package main

import (
	"os"
	"os/signal"

	"github.com/barnigator/sso/internal/app"

	"syscall"
)

func main() {
	application, err := app.Run()
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Stop()
}

// go run cmd/sso/main.go --config=./config/local.yaml
