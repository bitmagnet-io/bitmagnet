package tpdb

import (
	"errors"
	"git.sr.ht/~dragnel/go-tpdb"
	"github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config      Config
	Logger      *zap.SugaredLogger
	Search      search.Search
	Persistence persistence.Persistence
}

type Result struct {
	fx.Out
	TpdbClient *porndb.PorndbClient
	Client     Client
}

type Client interface {
	MovieClient
}

func New(p Params) (r Result, err error) {
	logger := p.Logger.Named("tpdb_client")

	if p.Config.ApiKey == defaultTpdbApiKey {
		logger.Warn("Metadataapi key not found")
		return Result{}, nil
	}

	logger.Infof("using : %s\n", p.Config.ApiKey)
	c := porndb.NewClient(p.Config.ApiKey)
	r.Client = &client{
		c: c,
		p: p.Persistence,
		s: p.Search,
	}
	r.TpdbClient = c
	return r, nil
}

type client struct {
	c *porndb.PorndbClient
	s search.Search
	p persistence.Persistence
}

const SourceTpdb = "tpdb"

var (
	ErrNotApi        = errors.New("no key found")
	ErrNotFound      = errors.New("not found")
	ErrUnknownSource = errors.New("unknown source")
)
