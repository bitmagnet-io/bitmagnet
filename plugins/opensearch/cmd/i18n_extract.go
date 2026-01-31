package main

import (
	"slices"

	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"github.com/bitmagnet-io/plugin-opensearch/config"
)

func main() {
	if err := write(); err != nil {
		panic(err)
	}
}

func write() error {
	return i18n.Write(
		"./i18n",
		slices.Concat(
			[]*i18n.Message{
				{
					ID:          "description",
					Description: "Description for the OpenSearch plugin",
					Other:       "OpenSearch indexer and search adapter",
				},
			},
			config.I18NMessages(),
		),
	)
}
