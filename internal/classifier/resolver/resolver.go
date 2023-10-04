package resolver

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	SubResolvers []SubResolver `group:"content_resolvers"`
	Dao          *dao.Query
	Logger       *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Resolver RootResolver
}

func New(p Params) (r Result, err error) {
	r.Resolver = &resolver{
		subResolvers: p.SubResolvers,
		dao:          p.Dao,
		logger:       p.Logger.Named("content_resolver"),
	}
	return
}

var (
	ErrNoMatch = errors.New("no match")
)

type Resolver interface {
	Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error)
}

type RootResolver interface {
	Resolver
	Persist(ctx context.Context, contents ...model.TorrentContent) error
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
	dao          *dao.Query
	logger       *zap.SugaredLogger
}
