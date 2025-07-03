package adapter

import (
	"context"
	"fmt"

	"github.com/autobrr/go-qbittorrent"
	"github.com/bitmagnet-io/bitmagnet/internal/client/model"
)

type qBitClient struct {
	CommonClient
}

func (c qBitClient) sendTo(ctx context.Context, content *content) error {
	sendTo, ok := c.config.GetSendTo(model.IDQBittorrent)
	if !ok {
		return model.ErrInvalidID
	}

	qb := qbittorrent.NewClient(qbittorrent.Config{
		Host:     fmt.Sprintf("http://%v:%v/", sendTo.Host, sendTo.Port),
		Username: sendTo.Username,
		Password: sendTo.Password,
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
