package processor

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Search     lazy.Lazy[search.Search]
	Classifier lazy.Lazy[classifier.Classifier]
	Dao        lazy.Lazy[*dao.Query]
	Logger     *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Processor lazy.Lazy[Processor]
}

func New(p Params) Result {
	return Result{
		Processor: lazy.New(func() (Processor, error) {
			s, err := p.Search.Get()
			if err != nil {
				return nil, err
			}
			d, err := p.Dao.Get()
			if err != nil {
				return nil, err
			}
			c, err := p.Classifier.Get()
			if err != nil {
				return nil, err
			}
			return processor{
				classifier: c,
				dao:        d,
				search:     s,
			}, nil
		}),
	}
}
