//go:build wasip1

package main

import (
	"context"
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/host/configurator"
	"github.com/bitmagnet-io/plugin-qbittorrent/config"
	"github.com/bitmagnet-io/plugin-qbittorrent/target"
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

	api.RegisterTorrentTarget(target.New(cfg))

	return nil
}
