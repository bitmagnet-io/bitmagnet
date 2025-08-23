package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/config/param"

type (
	Workflow    string
	Keywords    map[string][]string
	Extensions  map[string][]string
	FlagValues  map[string]any
	DeleteXXX   bool
	Concurrency int
)

var (
	ParamWorkflow = param.MustNew(
		param.WithDefault(Workflow("default")),
	)

	ParamKeywords = param.MustNew(
		param.WithMapstructure[Keywords](),
	)

	ParamExtensions = param.MustNew(
		param.WithMapstructure[Extensions](),
	)

	ParamFlags = param.MustNew(
		param.WithMapstructure[FlagValues](),
	)

	ParamDeleteXXX = param.MustNew[DeleteXXX]()

	ParamConcurrency = param.MustNew(
		param.WithDefault(Concurrency(100)),
		param.WithGreaterThan(Concurrency(0)),
	)
)

type Config struct {
	Workflow    Workflow
	Keywords    Keywords
	Extensions  Extensions
	Flags       FlagValues
	DeleteXXX   DeleteXXX
	Concurrency Concurrency
}

// func NewDefaultConfig() Config {
// 	return Config{
// 		Workflow:    "default",
// 		Concurrency: 100,
// 	}
// }
