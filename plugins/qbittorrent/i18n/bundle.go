package i18n

import (
	"embed"

	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
)

var (
	//go:embed *.yml
	embedFS embed.FS

	Bundle = i18n.MustNewBundle(fs.FromIOFS(embedFS))
)

func init() {
	Bundle.LanguageTags()
}
