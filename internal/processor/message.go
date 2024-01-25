package processor

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const MessageName = "process_torrent"

type ClassifyMode int

const (
	// ClassifyModeDefault will use any pre-existing content match as a hint
	// This is the most common use case and will only attempt to match previously unmatched torrents
	ClassifyModeDefault ClassifyMode = iota
	// ClassifyModeRematch will ignore any pre-existing classification and always classify from scratch
	// This is useful if there are errors in matches from an earlier version that need to be corrected
	ClassifyModeRematch
	// ClassifyModeSkipUnmatched will skip classification for previously unmatched torrents (that don't have any hint)
	// This is useful for eliminating expensive API calls while running the rest of the processing pipeline
	ClassifyModeSkipUnmatched
)

type MessageParams struct {
	ClassifyMode ClassifyMode
	InfoHashes   []protocol.ID
}
