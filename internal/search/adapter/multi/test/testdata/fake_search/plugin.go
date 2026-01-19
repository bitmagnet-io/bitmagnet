//go:build wasip1

package main

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/search"

	"github.com/bitmagnet-io/bitmagnet/proto/host/configurator"
)

type FakeTest struct {
	cfg configurator.Service
}

func main() {}

var plugin = &FakeTest{
	cfg: configurator.NewService(),
}

func init() {
	api.RegisterSearchAdapter(plugin)
}

func (t *FakeTest) SearchTorrentContent(ctx context.Context, params *search.Params) (*search.TorrentContentResult, error) {

	var c int32
	c = 99
	return &search.TorrentContentResult{
		TotalCount: &c,
	}, nil
}

func (t *FakeTest) SearchTorrentFiles(ctx context.Context, params *search.Params) (*search.TorrentFilesResult, error) {
	return nil, nil
}
