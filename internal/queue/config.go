package queue

import "github.com/bitmagnet-io/bitmagnet/internal/processor"

type Config struct {
	Concurrency int
	Queues      map[string]int
}

func NewDefaultConfig() Config {
	return Config{
		Concurrency: 10,
		Queues: map[string]int{
			processor.MessageName: 3,
		},
	}
}
