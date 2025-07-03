package adapter

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/client"
	clientmodel "github.com/bitmagnet-io/bitmagnet/internal/client/model"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type content = []search.TorrentContentResultItem

type clientWorker interface {
	AddInfoHashes(ctx context.Context, req client.AddInfoHashesRequest) error
	sendTo(ctx context.Context, content *content) error
}

type CommonClient struct {
	config *client.Config
	search search.Search
	client clientWorker
}

func New(cfg *client.Config, search search.Search) CommonClient {
	cc := CommonClient{
		config: cfg,
		search: search,
	}

	return cc
}

func (c CommonClient) downloadCategory(contentType model.ContentType) string {
	category := c.config.Categories[contentType]
	if category == "" {
		category = c.config.DefaultCategory
	}

	return category
}

func (c CommonClient) AddInfoHashes(ctx context.Context, req client.AddInfoHashesRequest) error {
	switch req.ClientID {
	case clientmodel.IDTransmission:
		c.client = transmissionClient{CommonClient: c}
	case clientmodel.IDQBittorrent:
		c.client = qBitClient{CommonClient: c}
	case clientmodel.IDNtfy:
		c.client = ntfy{CommonClient: c}
	default:
		return clientmodel.ErrInvalidID
	}

	options := []q.Option{
		q.Where(
			search.TorrentContentInfoHashCriteria(req.InfoHashes...),
		),
		search.TorrentContentCoreJoins(),
		search.HydrateTorrentContentContent(),
		search.HydrateTorrentContentTorrent(),
		q.Limit(uint(len(req.InfoHashes))),
	}

	sr, err := c.search.TorrentContent(ctx, options...)
	if err != nil {
		return err
	}

	return c.client.sendTo(ctx, &sr.Items)
}
