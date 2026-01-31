package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type (
	Workflow    string
	Keywords    map[string][]string
	Extensions  map[string][]string
	FlagValues  map[string]any
	DeleteXXX   bool
	Concurrency int
)

var (
	schemaMapStringStrings = json_schema.MustNew(
		json_schema.Typed(json_schema.TypeObject),
		json_schema.AdditionalPropertiesType(
			json_schema.MustNew(
				json_schema.Typed(json_schema.TypeArray),
				json_schema.Items(
					json_schema.MustNew(json_schema.Typed(json_schema.TypeString)),
				),
			),
		),
	)

	ParamWorkflow = param.MustNew(
		param.Description[Workflow]("The default classifier workflow"),
		param.Default(Workflow("default")),
	)

	ParamKeywords = param.MustNew(
		param.Description[Keywords](
			"A map of category names to keywords associated with different types of torrents",
		),
		param.Mapstructure[Keywords](),
		param.JSONSchema[Keywords](schemaMapStringStrings),
	)

	ParamExtensions = param.MustNew(
		param.Description[Extensions]("A map of file types to file extensions associated with them"),
		param.Mapstructure[Extensions](),
		param.JSONSchema[Extensions](schemaMapStringStrings),
	)

	ParamFlags = param.MustNew(
		param.Description[FlagValues]("A map of flag keys to flag values for configuring classifier workflows"),
		param.Mapstructure[FlagValues](),
	)

	ParamDeleteXXX = param.MustNew(
		param.Description[DeleteXXX](
			"A boolean flag indicating that XXX content should be deleted by the classifier",
		),
		param.Bool[DeleteXXX](),
	)

	ParamConcurrency = param.MustNew(
		param.Description[Concurrency]("Maximum number of classifications to run in parallel"),
		param.Int[Concurrency](),
		param.Default(Concurrency(100)),
		param.GreaterThan(Concurrency(0)),
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
