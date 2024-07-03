package client

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type AddInfoHashesRequest struct {
	ClientID   string
	InfoHashes []protocol.ID
}

type Client interface {
	AddInfoHashes(ctx context.Context, req AddInfoHashesRequest) error
}
