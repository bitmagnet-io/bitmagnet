package i18n

import "github.com/bitmagnet-io/bitmagnet/internal/slice"

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

type Provider func() []*Message

func NewProvider(messages ...*Message) Provider {
	return func() []*Message {
		return messages
	}
}

func Providers(providers ...Provider) Provider {
	return func() []*Message {
		return slice.FlatMap(providers, func(provider Provider) []*Message {
			return provider()
		})
	}
}
