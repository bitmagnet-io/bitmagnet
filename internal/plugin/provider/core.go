package provider

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	logging_console "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/console"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

var corePlugins = plugin.Plugins{
	config.Plugin,
	i18n.Plugin,
	logging.Plugin,
	logging_console.Plugin,
	postgres.Plugin,
}

func Core() plugin.Provider {
	return append(plugin.Plugins{}, corePlugins...)
}
