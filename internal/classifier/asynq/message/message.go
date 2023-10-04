package message

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const ClassifyTorrentTypename = "classify_torrent"

type ClassifyTorrentPayload struct {
	InfoHashes []model.Hash20
}
