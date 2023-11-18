package classifier

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/zap"
)

var (
	ErrNoMatch = errors.New("no match")
)

type Resolver interface {
	Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error)
}

type SubResolverConfig struct {
	Key      string
	Priority int
}

type SubResolver interface {
	Config() SubResolverConfig
	PreEnrich(content model.TorrentContent) (model.TorrentContent, error)
	Resolver
}

type resolver struct {
	subResolvers []SubResolver
	logger       *zap.SugaredLogger
}
