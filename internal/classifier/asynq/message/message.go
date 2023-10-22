package message

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const ClassifyTorrentTypename = "classify_torrent"

type ClassifyTorrentPayload struct {
	InfoHashes []protocol.ID
}
