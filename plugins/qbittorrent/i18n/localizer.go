package i18n

import (
	"embed"
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed *.yaml
	embedFS embed.FS

	bundle *i18n.Bundle
)

func init() {
	err := doInit()
	if err != nil {
		panic(err)
	}
}

func doInit() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	entries, err := fs.ReadDir(embedFS, ".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			bytes, err := embedFS.ReadFile(entry.Name())
			if err != nil {
				return err
			}

			_, err = bundle.ParseMessageFileBytes(bytes, entry.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewLocalizer(acceptLanguage []string) *Localizer {
	return &Localizer{i18n.NewLocalizer(bundle, acceptLanguage...)}
}

type Localizer struct {
	*i18n.Localizer
}

func (l *Localizer) Localize(messageID string) string {
	msg, err := l.Localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return messageID
	}
	return msg
}
