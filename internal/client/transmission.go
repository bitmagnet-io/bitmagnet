package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hekmon/transmissionrpc/v3"
)

type transmissionClient struct {
	commonClient
}

func (c transmissionClient) download(ctx context.Context, content *content, category string) error {

	endpoint, err := url.Parse(fmt.Sprintf("http://%v:%v/transmission/rpc", c.config.Transmission.Host, c.config.Transmission.Port))
	if err != nil {
		return err
	}
	tbt, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		return err
	}

	settings, err := tbt.SessionArgumentsGetAll(ctx)
	if err != nil {
		return err
	}

	dir := *settings.DownloadDir + "/" + category

	magnet := content.Torrent.MagnetUri()
	_, err = tbt.TorrentAdd(ctx, transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &dir,
	})
	return err

}
