//go:build dev

package bundle

import "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dev"

func init() {
	corePlugins = append(corePlugins, dev.Plugin)
}
