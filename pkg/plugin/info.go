package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/ref"

type PluginInfos []PluginInfo

type PluginInfo struct {
	Ref        ref.Ref
	Enabled    bool
	DependsOn  []ref.Ref
	RequiredBy []ref.Ref
}
