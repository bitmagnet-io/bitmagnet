//go:build dev || dev_all

package provider

import "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dev"

func init() {
	corePlugins = append(corePlugins, dev.Plugin)
}
