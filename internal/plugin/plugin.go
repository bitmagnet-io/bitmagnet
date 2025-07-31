package plugin

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"go.uber.org/fx"
)

type FXOptionProvider interface {
	FXOption() fx.Option
}

type Command func(app.Builder) cmd.Command

type CommandProvider interface {
	Commands() []Command
}

type Plugin interface {
	Ref() ref.Ref
	EnabledByDefault() bool
	Dependencies() []ref.Ref
	FXOptionProvider
	CommandProvider
}

func NewPlugin(
	ref ref.Ref,
	enabledByDefault bool,
	dependencies []ref.Ref,
	fxOption fx.Option,
	cliCommands []Command,
) Plugin {
	return &plugin{
		ref:              ref,
		enabledByDefault: enabledByDefault,
		dependencies:     dependencies,
		fxOption:         fxOption,
		cliCcommands:     cliCommands,
	}
}

type plugin struct {
	ref              ref.Ref
	enabledByDefault bool
	dependencies     []ref.Ref
	fxOption         fx.Option
	cliCcommands     []Command
}

func (p *plugin) Ref() ref.Ref {
	return p.ref
}

func (p *plugin) EnabledByDefault() bool {
	return p.enabledByDefault
}

func (p *plugin) Dependencies() []ref.Ref {
	return p.dependencies
}

func (p *plugin) FXOption() fx.Option {
	return fx.Module(
		p.Ref().String(),
		p.fxOption,
	)
}

func (p *plugin) Commands() []Command {
	if len(p.cliCcommands) == 0 {
		return nil
	}

	return p.cliCcommands
}
