package client

import (
	"context"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type AddInfoHashesRequest struct {
	ClientID   string
	InfoHashes []protocol.ID
}

type Client interface {
	AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error
}

type ServarrClient struct {
	Config Config
}

func (c *ServarrClient) AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error {
	strInfoHashes := make([]string, len(req.InfoHashes))
	for i, ih := range req.InfoHashes {
		strInfoHashes[i] = ih.String()
	}
	_, err := ServarrDownload(
		ctx,
		graphql.NewClient(c.Config.ArrServiceUrl+"/graphql", http.DefaultClient),
		strInfoHashes,
	)

	return err
}
