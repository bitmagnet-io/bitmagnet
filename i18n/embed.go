package i18n

import (
	"embed"

	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"github.com/spf13/afero"
)

var (
	//go:embed *.yaml
	embedFS embed.FS

	FS fs.FS = afero.Afero{Fs: afero.FromIOFS{FS: embedFS}}
)
