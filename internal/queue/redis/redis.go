package redis

import (
	r "github.com/redis/go-redis/v9"
)

type Client = r.Client

type Wrapper struct {
	Redis *r.Client
}

func (w Wrapper) MakeRedisClient() interface{} {
	return w.Redis
}
