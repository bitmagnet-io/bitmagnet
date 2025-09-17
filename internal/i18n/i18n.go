package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type (
	Bundle        = i18n.Bundle
	Localizer     = i18n.Localizer
	Message       = i18n.Message
	MessageOption func(*Message)
)

var (
	NewLocalizer = i18n.NewLocalizer
)

var languages = []language.Tag{
	language.Arabic,
	language.Catalan,
	language.German,
	language.English,
	language.Spanish,
	language.French,
	language.Hindi,
	language.Japanese,
	language.Portuguese,
	language.Russian,
	language.Turkish,
	language.Ukrainian,
	language.Chinese,
}
