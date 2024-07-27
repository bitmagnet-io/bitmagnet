package client

import (
	"context"
	"fmt"

	"github.com/autobrr/go-qbittorrent"
)

type qBitClient struct {
	commonClient
}

func (c qBitClient) download(ctx context.Context, content *content, category string) error {

	qb := qbittorrent.NewClient(qbittorrent.Config{
		Host:     fmt.Sprintf("http://%v:%v/", c.config.Qbittorrent.Host, c.config.Qbittorrent.Port),
		Username: c.config.Qbittorrent.Username,
		Password: c.config.Qbittorrent.Password,
	})

	err := qb.LoginCtx(ctx)
	if err != nil {
		return err
	}

	pref, err := qb.GetAppPreferencesCtx(ctx)
	if err != nil {
		return err
	}

	err = qb.AddTorrentFromUrlCtx(
		ctx,
		content.Torrent.MagnetUri(),
		map[string]string{
			"savepath": fmt.Sprintf("%v/%v", pref.SavePath, category),
			"category": category,
		},
	)

	return err

}
