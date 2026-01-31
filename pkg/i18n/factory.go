package i18n

import (
	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func NewBundle(inputFs fs.FS) (*Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	err := inputFs.Walk(".", func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			bytes, err := inputFs.ReadFile(path)
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

func MustNewBundle(inputFs fs.FS) *Bundle {
	bundle, err := NewBundle(inputFs)
	if err != nil {
		panic(err)
	}

	return bundle
}
