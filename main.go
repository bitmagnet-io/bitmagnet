package main

import (
	"github.com/bitmagnet-io/bitmagnet/internal/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app.New().Run()
}
