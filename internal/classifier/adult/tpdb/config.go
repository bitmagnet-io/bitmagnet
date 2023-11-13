package tpdb

type Config struct {
	ApiKey string
}

func NewDefaultConfig() Config {
	return Config{
		ApiKey: defaultTpdbApiKey,
	}
}

const (
	defaultTpdbApiKey = ""
)
