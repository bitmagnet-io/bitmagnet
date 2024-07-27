package client

import (
	"context"
	"fmt"

	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type AddInfoHashesRequest struct {
	ClientID   gen.ClientID
	InfoHashes []protocol.ID
}

type content = search.TorrentContentResultItem

type clientWorker interface {
	AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error
	download(ctx context.Context, content *content, category string) error
}

type commonClient struct {
	config *Config
	search search.Search
	client clientWorker
}

func New(cfg *Config, search search.Search) commonClient {
	cc := commonClient{
		config: cfg,
		search: search,
	}

	return cc

}

func (c commonClient) AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error {

	switch c.config.DownloadClient {
	case gen.ClientIDTransmission:
		c.client = transmissionClient{commonClient: c}
	case gen.ClientIDQBittorrent:
		c.client = qBitClient{commonClient: c}
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

	for _, cr := range sr.Items {
		category := c.config.Categories[cr.Content.Type]
		if category == "" {
			category = c.config.DefaultCategory
		}
		err = c.client.download(ctx, &cr, category)
		if err != nil {
			return err
		}

	}

	return nil
}
