package main

import (
	"os"

	"github.com/bitmagnet-io/bitmagnet/internal/app"
	"github.com/bitmagnet-io/bitmagnet/internal/env"
)

func main() {
	exitCode, _ := app.Run(env.NewDefault())

	os.Exit(exitCode)
}
