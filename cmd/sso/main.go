package main

import (
	"fmt"
	"sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Print(cfg)
}

// go run cmd/sso/main.go --config=./config/local.yaml
