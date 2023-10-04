// this package provides a client for itorrents.org.
// their API is behind CloudFlare and any requests need heavily rate limiting.
// I haven't found a good use for it yet but it may be used for ad-hoc fetching of torrent metadata at some point.

package itorrents

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpclient/httplogger"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpclient/httpratelimiter"
	"github.com/bitmagnet-io/bitmagnet/internal/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

type Params struct {
	fx.In
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Client Client
}

type Client interface {
	Get(ctx context.Context, hash20 model.Hash20) (metainfo.TorrentFile, error)
}

func New(p Params) Result {
	return Result{
		Client: client{
			http: http.Client{
				Transport: httpratelimiter.NewDecorator(
					time.Second,
					1,
				)(httplogger.NewDecorator(
					p.Logger.Named("itorrents_client"),
				)(http.DefaultTransport)),
			},
		},
	}
}

type client struct {
	http http.Client
}

func (c client) Get(ctx context.Context, hash20 model.Hash20) (metainfo.TorrentFile, error) {
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, "https://itorrents.org/torrent/"+strings.ToUpper(hash20.String())+".torrent", nil)
	if reqErr != nil {
		return metainfo.TorrentFile{}, fmt.Errorf("failed to create request: %w", reqErr)
	}
	res, resErr := c.http.Do(req)
	if resErr != nil {
		return metainfo.TorrentFile{}, fmt.Errorf("failed to request torrent file: %w", reqErr)
	}
	if res.StatusCode != http.StatusOK {
		return metainfo.TorrentFile{}, fmt.Errorf("response status was not OK requesting torrent file: %s", res.Status)
	}
	bytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return metainfo.TorrentFile{}, fmt.Errorf("failed to read torrent file: %w", readErr)
	}
	return metainfo.ReadTorrentFileBytes(bytes)
}
