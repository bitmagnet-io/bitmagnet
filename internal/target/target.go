package target

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_forms"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type Target[T any] interface {
	Ref() ref.Ref
	Name() string
	DataSchama(ctx context.Context) (*json_schema.JSONSchema, error)
	UISchema(ctx context.Context, acceptLanguage []string) (*json_forms.UISchema, error)
	Send(ctx context.Context, items []T, data json_schema.JSONValue) (*json_schema.JSONValue, error)
}

type TorrentContentTarget = Target[model.TorrentContent]
