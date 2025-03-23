package gormcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	DB lazy.Lazy[*gorm.DB]
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (r Result, err error) {
	r.Command = &cli.Command{
		Name: "gorm",
		Subcommands: []*cli.Command{
			{
				Name: "gen",
				Action: func(*cli.Context) error {
					db, err := p.DB.Get()
					if err != nil {
						return err
					}
					g := gen.BuildGenerator(db)
					g.Execute()
					return nil
				},
			},
		},
	}
	return
}
