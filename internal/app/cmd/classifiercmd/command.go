package classifiercmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

type Params struct {
	fx.In
	WorkflowSource lazy.Lazy[classifier.WorkflowSource]
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	return Result{Command: &cli.Command{
		Name: "classifier",
		Subcommands: []*cli.Command{
			{
				Name:  "show",
				Usage: "Show the classifier workflow source",
				Action: func(ctx *cli.Context) error {
					src, srcErr := p.WorkflowSource.Get()
					if srcErr != nil {
						return srcErr
					}
					y, yErr := yaml.Marshal(src)
					if yErr != nil {
						return yErr
					}
					println(string(y))
					return nil
				},
			},
		},
	}}, nil
}
