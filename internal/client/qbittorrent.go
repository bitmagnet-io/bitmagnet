package client

import (
	"context"
	"fmt"

	"github.com/autobrr/go-qbittorrent"
)

type qBitClient struct {
	CommonClient
}

func (c qBitClient) download(ctx context.Context, content *content) error {
	qb := qbittorrent.NewClient(qbittorrent.Config{
		Host:     fmt.Sprintf("http://%v:%v/", c.config.Qbittorrent.Host, c.config.Qbittorrent.Port),
		Username: c.config.Qbittorrent.Username,
		Password: c.config.Qbittorrent.Password,
		Timeout:  1,
	})

	err := qb.LoginCtx(ctx)
	if err != nil {
		return err
	}

	pref, err := qb.GetAppPreferencesCtx(ctx)
	if err != nil {
		return err
	}

	for _, item := range *content {
		category := c.downloadCategory(item.Content.Type)

		err = qb.AddTorrentFromUrlCtx(
			ctx,
			item.Torrent.MagnetURI(),
			map[string]string{
				"savepath": fmt.Sprintf("%v/%v", pref.SavePath, category),
				"category": category,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
