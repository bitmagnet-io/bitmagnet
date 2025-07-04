package configcmd

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/app"
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Command = &cli.Command{
	Name: "config",
	Commands: []*cli.Command{
		{
			Name:  "show",
			Usage: "Shows all available configuration parameters and resolved values",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fxApp := app.New(
					ctx,
					cmd.Writer,
					fx.Invoke(func(
						logger *zap.SugaredLogger,
						cfg config.ResolvedConfig,
					) error {
						logger.Debugw("resolved config", "cfg", cfg)
						tw := table.NewWriter()
						tw.SetOutputMirror(cmd.Writer)
						tw.AppendHeader(table.Row{"path", "Type", "Value", "Default", "From"})
						for _, node := range cfg.Nodes() {
							appendRows(tw, node)
						}
						tw.Render()
						return nil
					}),
				)

				fxApp.Run()

				return fxApp.Err()
			},
		},
	},
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
