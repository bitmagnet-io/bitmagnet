package client

import (
	"context"
	"fmt"

	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type AddInfoHashesRequest struct {
	ClientID   gen.ClientID
	InfoHashes []protocol.ID
}

type content = []search.TorrentContentResultItem

type clientWorker interface {
	AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error
	download(ctx context.Context, content *content) error
}

type CommonClient struct {
	config *Config
	search search.Search
	client clientWorker
}

func New(cfg *Config, search search.Search) CommonClient {
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

func (c CommonClient) AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error {
	switch c.config.DownloadClient {
	case gen.ClientIDTransmission:
		c.client = transmissionClient{CommonClient: c}
	case gen.ClientIDQBittorrent:
		c.client = qBitClient{CommonClient: c}
	default:
		return fmt.Errorf("not implemented %s", c.config.DownloadClient)
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

	return c.client.download(ctx, &sr.Items)
}
