package config

// import (
// 	"context"

// 	"github.com/bitmagnet-io/bitmagnet/internal/env"
// 	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
// 	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
// 	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
// 	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
// 	"github.com/charmbracelet/bubbles/table"
// 	"github.com/charmbracelet/lipgloss"
// 	"go.uber.org/fx"
// )

// func NewConfigCommand() plugin.Command {
// 	return func(appBuilder app.Builder) cmd.Command {
// 		return &ConfigCommand{
// 			App: cmd.NewApp[ConfigDeps](appBuilder),
// 		}
// 	}
// }

// type ConfigDeps struct {
// 	fx.In
// 	ResolvedConfig config.ResolvedConfig
// }

// type ConfigCommand struct {
// 	cmd.Cmd
// 	cmd.App[ConfigDeps]
// }

// func (cfgCmd *ConfigCommand) Subcommands() []cmd.Command {
// 	return []cmd.Command{
// 		&ShowCommand{
// 			Cmd: cfgCmd.Cmd,
// 			App: cfgCmd.App,
// 		},
// 	}
// }

// type ShowCommand struct {
// 	cmd.Cmd
// 	cmd.App[ConfigDeps]
// }

// func (cmd *ShowCommand) Run(env env.Env) error {
// 	return cmd.NewRunner(func(deps ConfigDeps) runner.Runner {
// 		return runner.SimpleRunner(func(context.Context) error {
// 			style := lipgloss.NewStyle().
// 				Bold(true).
// 				// Italic(true).
// 				// Faint(true).
// 				Blink(true).
// 				// Strikethrough(true).
// 				// Underline(true).
// 				// Reverse(true).
// 				Foreground(lipgloss.Color("#4355b9"))
// 			env.Write([]byte(style.Render("hello") + "\n"))
// 			_ = lipgloss.Color("#4355b9")
// 			m := table.New(
// 				table.WithColumns([]table.Column{
// 					{Title: "Path", Width: 20},
// 					{Title: "Type", Width: 10},
// 					{Title: "Value", Width: 30},
// 					{Title: "Default", Width: 30},
// 					{Title: "From", Width: 20},
// 				}),
// 				table.WithRows([]table.Row{
// 					{"path1", "string", "value1", "default1", "source1"},
// 				}),
// 			)
// 			env.Write([]byte(m.View()))
// 			return nil
// 		})
// 	})(env)
// }
