package cache

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/lru"
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/zap"
)

func New(
	ttl *atomic.Value[TTL],
	maxItems *atomic.Value[MaxItems],
	logger *zap.Logger,
) caches.Cacher {
	currentTTL, currentMaxItems := ttl.Get(), maxItems.Get()

	cacher := &inMemoryCacher{
		lru: lru.New[string, *caches.Query](
			int(currentMaxItems),
			nil,
			time.Duration(currentTTL),
		),
		logger: logger,
	}

	ttl.Subscribe(func(ttl TTL) {
		if ttl != currentTTL {
			currentTTL = ttl
			cacher.lru.SetTTL(time.Duration(currentTTL))
		}
	})

	maxItems.Subscribe(func(maxItems MaxItems) {
		if maxItems != currentMaxItems {
			currentMaxItems = maxItems
			cacher.lru.Resize(int(currentMaxItems))
		}
	})

	return cacher
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

type modeKeyType string

// The modeKey context value specifies the caching Mode for a particular query.
// I don't really like storing this in the context, but it's the simplest way for now.
const modeKey modeKeyType = "gorm_cache_mode"

type inMemoryCacher struct {
	lru    *lru.LRU[string, *caches.Query]
	logger *zap.Logger
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

	c.logger.Debug("cache hit", zap.String("key", key))

	return val
}

func (c *inMemoryCacher) Store(ctx context.Context, key string, val *caches.Query) error {
	m := cacheModeFromContext(ctx)
	if m == ModeCached || m == ModeWarm {
		c.lru.Add(key, val)
	}

	return nil
}

func ContextWithCacheMode(ctx context.Context, mode Mode) context.Context {
	return context.WithValue(ctx, modeKey, mode)
}

func cacheModeFromContext(ctx context.Context) Mode {
	ctxValue := ctx.Value(modeKey)

	m, isOk := ctxValue.(Mode)
	if !isOk {
		return ModeNoCache
	}

	return m
}
