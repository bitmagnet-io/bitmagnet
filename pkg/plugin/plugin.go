package plugin

import (
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"go.uber.org/fx"
)

type ConfigProvider interface {
	ConfigParams() []config_registry.Param
}

type ErrorProvider interface {
	Errors() ref.Map[error]
}

type I18nProvider interface {
	I18nMessages() []*i18n.Message
}

type Command func(app.Builder) cmd.Command

type CommandProvider interface {
	Commands() []Command
}

type FXOptionProvider interface {
	FXOption() fx.Option
}

type Plugin interface {
	Ref() ref.Ref
	ActivationRef() ref.Nullable
	Dependencies() []ref.Ref
	ConfigProvider
	ErrorProvider
	I18nProvider
	LocalizedContentProvider
	CommandProvider
	FXOptionProvider
}
