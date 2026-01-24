//go:build wasip1

package target

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/bitmagnet/proto/common/model"
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/plugin-qbittorrent/config"
)

type target struct {
	config     config.Config
	httpClient http_client.Service
}

func New(config config.Config) api.TorrentTarget {
	return &target{
		config:     config,
		httpClient: http_client.NewService(),
	}
}

func (t *target) request(ctx context.Context, req *http.Request, result any) error {
	res, err := t.httpClient.Request(ctx, req)
	if err != nil {
		return err
	}

	if res.Status == 401 || res.Status == 403 {
		return errUnauthorized
	}

	if res.Status < 200 || res.Status >= 300 {
		return fmt.Errorf("request failed: status code: %d", res.Status)
	}

	if string(res.Body) == "Fails." {
		return errors.New("request failed")
	}

	if result != nil {
		err = json.Unmarshal(res.Body, result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

var errUnauthorized = errors.New("unauthorized")

func (t *target) login(ctx context.Context) error {
	err := t.request(ctx, &http.Request{
		Url:    t.config.URL + "/api/v2/auth/login",
		Method: http.Method_post,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Body: []byte(fmt.Sprintf("username=%s&password=%s", t.config.Username, t.config.Password)),
	}, nil)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	return nil
}

func (t *target) retryWithLogin(ctx context.Context, fn func() error) error {
	err := fn()
	if errors.Is(err, errUnauthorized) {
		err = t.login(ctx)
		if err == nil {
			err = fn()
		}
	}

	return err
}

func createMagnetLink(torrent *model.Torrent) string {
	return "magnet:?xt=urn:btih:" + torrent.InfoHash + "&dn=" + torrent.Name
}
