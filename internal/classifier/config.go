package classifier

type Config struct {
	DefaultWorkflow string
	Keywords        map[string][]string
	Extensions      map[string][]string
	Flags           map[string]any
}

func NewDefaultConfig() Config {
	return Config{
		DefaultWorkflow: "default",
	}
}
