package httpserver

type Config struct {
	LocalAddress string
	GinMode      string
	Options      []string
}

func NewDefaultConfig() Config {
	return Config{
		LocalAddress: ":3333",
		GinMode:      "release",
		Options:      []string{"*"},
	}
}
