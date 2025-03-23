package configcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config config.ResolvedConfig
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	cmd := &cli.Command{
		Name: "config",
		Subcommands: []*cli.Command{
			{
				Name:  "show",
				Usage: "Shows all available configuration parameters and resolved values",
				Action: func(ctx *cli.Context) error {
					p.Logger.Debugw("resolved config", "cfg", p.Config)
					tw := table.NewWriter()
					tw.SetOutputMirror(ctx.App.Writer)
					tw.AppendHeader(table.Row{"path", "Type", "Value", "Default", "From"})
					for _, node := range p.Config.Nodes() {
						appendRows(tw, node)
					}
					tw.Render()
					return nil
				},
			},
		},
	}

	return Result{Command: cmd}, nil
}

func appendRows(tw table.Writer, node config.ResolvedNode) {
	if node.IsStruct {
		tw.AppendSeparator()
		tw.AppendRow(table.Row{node.PathString + ":", node.Type})
		tw.AppendSeparator()

		for _, child := range node.Children() {
			appendRows(tw, child)
		}

		tw.AppendSeparator()
	} else {
		from := node.ResolverKey
		if from == "" {
			from = "default"
		}

		tw.AppendRow(table.Row{node.PathString, node.Type, node.ValueLabel, node.DefaultLabel, from})
	}
}
