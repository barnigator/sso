package main

import (
	"os"
	"os/signal"
	"sso/internal/app"
	"syscall"
)

func main() {

	application, err := app.Run()
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	application.GRPCServer.Stop(stopSign)
}

// go run cmd/sso/main.go --config=./config/local.yaml
