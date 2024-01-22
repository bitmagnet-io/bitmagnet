package processor

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const MessageName = "process_torrent"

type MessageParams struct {
	Rematch    bool
	InfoHashes []protocol.ID
}
