package adapter

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/client/model"
	"github.com/go-resty/resty/v2"
)

type ntfy struct {
	CommonClient
}

func (n ntfy) sendTo(ctx context.Context, content *content) error {
	sendTo, ok := n.config.GetSendTo(model.IDNtfy)
	if !ok {
		return model.ErrInvalidID
	}

	topic, exists := n.config.Categories["ntfy"]
	if !exists {
		topic = "magnet"
	}

	r := resty.New().R().SetContext(ctx)

	for _, item := range *content {
		resp, err := r.SetBody(item.Torrent.MagnetURI()).
			Post(fmt.Sprintf("http://%v:%v/%v", sendTo.Host, sendTo.Port, topic))
		if err != nil {
			return err
		} else if resp.IsError() {
			return fmt.Errorf("[%v] %v", resp.StatusCode(), resp.Request.URL)
		}
	}

	return nil
}
