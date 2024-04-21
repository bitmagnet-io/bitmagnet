package classifier

type Config struct {
	Workflow   string
	Keywords   map[string][]string
	Extensions map[string][]string
	Flags      map[string]any
	DeleteXxx  bool
}

func NewDefaultConfig() Config {
	return Config{
		Workflow: "default",
	}
}
