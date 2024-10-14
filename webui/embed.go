package webui

import "embed"

//go:embed dist/bitmagnet/browser/*
var staticFS embed.FS

func StaticFS() embed.FS {
	return staticFS
}
