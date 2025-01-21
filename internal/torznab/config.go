package torznab

type Config struct {
	Hostname *string
	Profiles map[string]map[string]string
}

func NewDefaultConfig() Config {
	return Config{}
}
