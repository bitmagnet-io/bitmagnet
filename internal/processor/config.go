package processor

type Config struct {
	Concurrency uint
}

func NewDefaultConfig() Config {
	return Config{
		Concurrency: 1,
	}
}
