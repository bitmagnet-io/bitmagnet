//go:build wasip1

package target

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
)

func (t *target) Send(ctx context.Context, sendTorrentsParams *api.SendTorrentsParams) (*api.JSONPayload, error) {
	var data Data
	err := json.Unmarshal(sendTorrentsParams.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	field, err := w.CreateFormField("urls")
	if err != nil {
		return nil, err
	}

	for _, torrent := range sendTorrentsParams.Torrents {
		_, err := io.WriteString(field, createMagnetLink(torrent.Torrent)+"\n")
		if err != nil {
			return nil, err
		}
	}

	if data.Category != "" {
		err = w.WriteField("category", data.Category)
		if err != nil {
			return nil, err
		}
	}

	if data.Stopped {
		err = w.WriteField("stopped", "true")
		if err != nil {
			return nil, err
		}
		err = w.WriteField("paused", "true")
		if err != nil {
			return nil, err
		}
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	bytes := buf.Bytes()

	err = t.retryWithLogin(ctx, func() error {
		return t.request(ctx, &http.Request{
			Url:    t.config.URL + "/api/v2/torrents/add",
			Method: http.Method_post,
			Headers: map[string]string{
				"Content-Type":   w.FormDataContentType(),
				"Content-Length": fmt.Sprintf("%d", len(bytes)),
			},
			Body: bytes,
		}, nil)
	})

	if err != nil {
		return nil, err
	}

	return &api.JSONPayload{}, nil
}
