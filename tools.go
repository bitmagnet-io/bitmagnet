//go:build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/abice/go-enum"
	_ "github.com/vektra/mockery/v2"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
