package client

type Config struct {
	ArrServiceUrl string
}

func NewDefaultConfig() Config {
	return Config{
		ArrServiceUrl: "http://localhost:3335",
	}
}
