package main

import (
	"os"

	"github.com/bitmagnet-io/bitmagnet/pkg/app"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
)

func main() {
	exitCode, _ := app.New().Run(env.NewDefault())

	os.Exit(exitCode)
}
