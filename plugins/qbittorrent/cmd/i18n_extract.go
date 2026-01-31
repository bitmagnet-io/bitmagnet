package main

import (
	"slices"

	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"github.com/bitmagnet-io/plugin-qbittorrent/config"
	"github.com/bitmagnet-io/plugin-qbittorrent/target"
)

func main() {
	if err := write(); err != nil {
		panic(err)
	}
}

func write() error {
	return i18n.Write("./i18n", slices.Concat(
		[]*i18n.Message{
			{
				ID:          "description",
				Description: "Description for the qBittorrent plugin",
				Other:       "qBittorrent target",
			},
		},
		config.I18NMessages(),
		target.I18NMessages(),
	))
}
