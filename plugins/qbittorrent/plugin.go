//go:build wasip1

package main

import (
	"context"

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
	cfg, err := configurator.Resolve[config.Config](ctx)
	if err != nil {
		return err
	}

	api.RegisterTorrentTarget(target.New(cfg))

	return nil
}
