package bundle

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	database_postgres "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	logging_console "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/console"
)

var corePlugins = []plugin.Plugin{
	config.Plugin,
	database_postgres.Plugin,
	// errors.Plugin,
	i18n.Plugin,
	logging.Plugin,
	logging_console.Plugin,
}

func Core() (Bundle, error) {
	return New(core.Ref, corePlugins...)
}
