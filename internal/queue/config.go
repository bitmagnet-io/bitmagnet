package queue

type Config struct {
	Concurrency int
	Queues      map[string]int
}

func NewDefaultConfig() Config {
	return Config{
		Concurrency: 10,
		Queues: map[string]int{
			"classify_torrent": 3,
		},
	}
}
