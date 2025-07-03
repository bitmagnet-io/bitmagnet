package client

import (
	"github.com/bitmagnet-io/bitmagnet/internal/client/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type AddInfoHashesRequest struct {
	ClientID   model.ID
	InfoHashes []protocol.ID
}
