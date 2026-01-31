package i18n

import "github.com/bitmagnet-io/bitmagnet/pkg/i18n"

type (
	Localizer struct {
		*i18n.Localizer
	}

	Message = i18n.Message
)

func NewLocalizer(acceptLanguage []string) *Localizer {
	return &Localizer{i18n.NewLocalizer(Bundle, acceptLanguage...)}
}

func (l *Localizer) Localize(messageID string) string {
	if msg, _ := l.Localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID}); msg != "" {
		return msg
	}

	return messageID
}
