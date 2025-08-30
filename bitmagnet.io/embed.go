package docs

import "embed"

//go:embed _site/*
var staticFS embed.FS

func StaticFS() embed.FS {
	return staticFS
}
