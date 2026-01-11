//go:build wasip1

package main

import (
	"context"
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/host/configurator"
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/plugin-opensearch/client"
	"github.com/bitmagnet-io/plugin-opensearch/config"
	"github.com/bitmagnet-io/plugin-opensearch/indexer"
	"github.com/bitmagnet-io/plugin-opensearch/search"
)

func main() {}

func init() {
	err := doInit(context.Background())
	if err != nil {
		panic(err)
	}
}

func doInit(ctx context.Context) error {
	rawConfig, err := configurator.NewService().GetConfig(ctx, nil)
	if err != nil {
		return err
	}

	var cfg config.Config
	err = json.Unmarshal([]byte(rawConfig.Json), &cfg)
	if err != nil {
		return err
	}

	client, err := client.New(http_client.NewService(), cfg)
	if err != nil {
		return err
	}

	api.RegisterIndexer(indexer.New(client, cfg.IndexPrefix))
	api.RegisterSearchAdapter(search.New(client, cfg.IndexPrefix))

	return nil
}
