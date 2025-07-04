package main

import (
	"context"
	"os"

	"github.com/bitmagnet-io/bitmagnet/internal/cli"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := cli.App.Run(context.Background(), os.Args)

	exitCode := 0

	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
