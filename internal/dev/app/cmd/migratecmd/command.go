package migratecmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrations"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Migrator lazy.Lazy[migrations.Migrator]
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (r Result, err error) {
	r.Command = &cli.Command{
		Name: "migrate",
		Subcommands: []*cli.Command{
			{
				Name: "up",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "version",
						Value: 0,
					},
				},
				Action: func(ctx *cli.Context) error {
					m, err := p.Migrator.Get()
					if err != nil {
						return err
					}
					version := ctx.Int64("version")
					if version == 0 {
						return m.Up(ctx.Context)
					} else {
						return m.UpTo(ctx.Context, version)
					}
				},
			},
			{
				Name: "down",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "version",
						Value: 0,
					},
				},
				Action: func(ctx *cli.Context) error {
					m, err := p.Migrator.Get()
					if err != nil {
						return err
					}
					version := ctx.Int64("version")
					if version == 0 {
						return m.Down(ctx.Context)
					} else {
						return m.DownTo(ctx.Context, version)
					}
				},
			},
		},
	}
	return
}
