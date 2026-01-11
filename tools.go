//go:build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/abice/go-enum"
	_ "github.com/go-delve/delve/cmd/dlv"
	_ "github.com/nicksnyder/go-i18n/v2/goi18n"
	_ "github.com/vektra/mockery/v2"
	_ "golang.org/x/tools/gopls"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "mvdan.cc/gofumpt"
)
