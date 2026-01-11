package main

import (
	"flag"

	"github.com/bitmagnet-io/bitmagnet/internal/wasm/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet

	protogen.Options{ParamFunc: flags.Set}.Run(gen.GeneratePlugin)
}
