package plugin

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type LocalizedContent struct {
	Ref          ref.Ref                `json:"ref"`
	Description  string                 `json:"description,omitempty"`
	ConfigParams []LocalizedConfigParam `json:"configParams,omitempty"`
}

type LocalizedConfigParam struct {
	Ref         ref.Ref `json:"ref"`
	Description string  `json:"description,omitempty"`
}

type LocalizedContentProvider interface {
	LocalizedContent(ctx context.Context, acceptLanguage ...string) (LocalizedContent, error)
}
