package main

import (
	"github.com/mqdvi-dp/go-common/config"
	"github.com/mqdvi-dp/go-common/example/example-service/cmd"
)

const serviceName = "example"

func main() {
	cfg := config.New(serviceName)
	defer cfg.Exit()

	srv := cmd.Serve(cfg)
	srv.Run()
}
