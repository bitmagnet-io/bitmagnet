package i18n

import (
	"io/fs"

	embed_i18n "github.com/bitmagnet-io/bitmagnet/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func NewBundle() (*Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	err := embed_i18n.FS.Walk(".", func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			bytes, err := embed_i18n.FS.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = bundle.ParseMessageFileBytes(bytes, path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return bundle, err
}
