package workercmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/worker"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Workers worker.Registry
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	cmd := &cli.Command{
		Name: "worker",
		Subcommands: []*cli.Command{
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Value: false,
					},
					&cli.StringSliceFlag{
						Name: "keys",
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.Bool("all") {
						p.Workers.EnableAll()
					} else {
						if err := p.Workers.Enable(ctx.StringSlice("keys")...); err != nil {
							return err
						}
					}
					if err := p.Workers.Start(ctx.Context); err != nil {
						return err
					}
					<-ctx.Context.Done()
					return nil
				},
				After: func(context *cli.Context) error {
					return p.Workers.Stop(context.Context)
				},
			},
			{
				Name: "list",
				Action: func(*cli.Context) error {
					for _, w := range p.Workers.Workers() {
						println(w.Key())
					}
					return nil
				},
			},
		},
	}

	return Result{Command: cmd}, nil
}
