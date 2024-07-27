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

type ServicesClient struct {
	Config Config
}

func (c *ServicesClient) AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error {
	_, err := ServicesDownload(
		ctx,
		graphql.NewClient(c.Config.ArrServiceUrl+"/graphql", http.DefaultClient),
		req.InfoHashes,
		ClientID(req.ClientID),
	)

	return err
}
