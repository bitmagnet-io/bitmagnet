package redis

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	r "github.com/redis/go-redis/v9"
)

type Client = r.Client

type Wrapper struct {
	Redis lazy.Lazy[*r.Client]
}

func (w Wrapper) MakeRedisClient() interface{} {
	redis, err := w.Redis.Get()
	if err != nil {
		return err
	}
	return redis
}
