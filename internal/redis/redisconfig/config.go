package redisconfig

import (
	redis "github.com/redis/go-redis/v9"
	"time"
)

func NewDefaultConfig() Config {
	return Config{
		Network:      "tcp",
		Addr:         "localhost:6379",
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	}
}

type Config struct {
	// Network type to use, either tcp or unix.
	// Default is tcp.
	Network string

	// Redis server address in "host:port" format.
	Addr string

	// Username to authenticate the current connection when Redis ACLs are used.
	// See: https://redis.io/commands/auth.
	Username string

	// Password to authenticate the current connection.
	// See: https://redis.io/commands/auth.
	Password string

	// Redis DB to select after connecting to a server.
	// See: https://redis.io/commands/select.
	DB int

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration

	// Timeout for socket reads.
	// If timeout is reached, read commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration

	// Timeout for socket writes.
	// If timeout is reached, write commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is ReadTimout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
}

func (c Config) RedisClientOptions() *redis.Options {
	return &redis.Options{
		Network:      c.Network,
		Addr:         c.Addr,
		Username:     c.Username,
		Password:     c.Password,
		DB:           c.DB,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
	}
}
