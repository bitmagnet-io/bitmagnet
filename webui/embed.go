package webui

import "embed"

//go:embed dist/*
var staticFS embed.FS

func StaticFS() embed.FS {
	return staticFS
}
