package cache

import (
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/fx"
)

type PluginParams struct {
	fx.In
	Config Config
	Cacher caches.Cacher
}

type PluginResult struct {
	fx.Out
	Plugin *caches.Caches
}

func NewPlugin(p PluginParams) PluginResult {
	var cacher caches.Cacher
	if p.Config.CacheEnabled {
		cacher = p.Cacher
	}

	return PluginResult{
		Plugin: &caches.Caches{Conf: &caches.Config{
			Easer:  p.Config.EaserEnabled,
			Cacher: cacher,
		}},
	}
}
