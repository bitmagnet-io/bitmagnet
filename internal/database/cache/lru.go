package cache

import (
	"context"

	"github.com/hashicorp/golang-lru/v2/expirable"
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config Config
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Cacher caches.Cacher
}

func NewInMemoryCacher(p Params) Result {
	return Result{
		Cacher: &inMemoryCacher{
			lru: expirable.NewLRU[string, *caches.Query](
				int(p.Config.MaxKeys),
				nil,
				p.Config.TTL,
			),
			logger: p.Logger.Named("gorm_cache"),
		},
	}
}

type Mode int

const (
	// ModeNoCache the query will not be satisfied from the cache,
	// and any existing cache entry will be removed to avoid stale results in future (default)
	ModeNoCache Mode = iota
	// ModeCached the query will be satisfied from the cache if possible,
	// otherwise the result will be stored in the cache
	ModeCached
	// ModeWarm the query will not be satisfied from the cache,
	// but the result will be stored in the cache for future queries using ModeCached
	ModeWarm
)

type modeKey string

// The ModeKey context value specifies the caching Mode for a particular query.
// I don't really like storing this in the context, but it's the simplest way for now.
const ModeKey modeKey = "gorm_cache_mode"

type inMemoryCacher struct {
	lru    *expirable.LRU[string, *caches.Query]
	logger *zap.SugaredLogger
}

func (c *inMemoryCacher) Get(ctx context.Context, key string) *caches.Query {
	m := cacheModeFromContext(ctx)
	if m == ModeNoCache || m == ModeWarm {
		return nil
	}

	val, ok := c.lru.Get(key)
	if !ok {
		return nil
	}

	c.logger.Debugw("cache hit", "key", key)

	return val
}

func (c *inMemoryCacher) Store(ctx context.Context, key string, val *caches.Query) error {
	m := cacheModeFromContext(ctx)
	if m == ModeCached || m == ModeWarm {
		c.lru.Add(key, val)
	}

	return nil
}

func cacheModeFromContext(ctx context.Context) Mode {
	ctxValue := ctx.Value(ModeKey)

	m, isOk := ctxValue.(Mode)
	if !isOk {
		return ModeNoCache
	}

	return m
}
