package processor

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
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
)

type MessageParams struct {
	ClassifyMode       ClassifyMode     `json:"ClassifyMode,omitempty"`
	ClassifierWorkflow string           `json:"ClassifierWorkflow,omitempty"`
	ClassifierFlags    classifier.Flags `json:"ClassifierFlags,omitempty"`
	InfoHashes         []protocol.ID    `json:"InfoHashes"`
}

func NewQueueJob(msg MessageParams, options ...model.QueueJobOption) (model.QueueJob, error) {
	return model.NewQueueJob(
		MessageName,
		msg,
		append([]model.QueueJobOption{model.QueueJobMaxRetries(2)}, options...)...,
	)
}
