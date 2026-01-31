package i18n

import (
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"golang.org/x/text/language"
)

func NewMessage(id, description string, options ...MessageOption) *Message {
	msg := &Message{
		ID:          id,
		Description: description,
	}

	for _, option := range options {
		option(msg)
	}

	return msg
}

func WithOther(other string) MessageOption {
	return func(msg *Message) {
		msg.Other = other
	}
}

type MessageProvider func() []*Message

func NewMessageProvider(messages ...*Message) MessageProvider {
	return func() []*Message {
		return messages
	}
}

func MessageProviders(providers ...MessageProvider) MessageProvider {
	return func() []*Message {
		return slice.FlatMap(providers, func(provider MessageProvider) []*Message {
			return provider()
		})
	}
}

type TranslationProvider func() map[language.Tag][]*Message
