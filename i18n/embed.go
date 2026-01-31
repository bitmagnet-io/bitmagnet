package i18n

import (
	"embed"

	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"github.com/spf13/afero"
)

var (
	//go:embed *.yml
	embedFS embed.FS

	Bundle = i18n.MustNewBundle(afero.Afero{Fs: afero.FromIOFS{FS: embedFS}})
)
