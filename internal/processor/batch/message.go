package batch

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const MessageName = "process_torrent_batch"

type MessageParams struct {
	InfoHashGreaterThan protocol.ID
	ClassifyMode        processor.ClassifyMode
	ClassifierWorkflow  string
	ClassifierFlags     classifier.Flags
	ChunkSize           uint
	BatchSize           uint
	ContentTypes        []model.NullContentType
	Orphans             bool
}

func (p MessageParams) ApisDisabled() bool {
	if p.ClassifierFlags == nil {
		return false
	}
	enabledAny, ok := p.ClassifierFlags["apis_enabled"]
	if !ok {
		return false
	}
	enabled, ok := enabledAny.(bool)
	return ok && !enabled
}

func NewQueueJob(msg MessageParams, options ...model.QueueJobOption) (model.QueueJob, error) {
	return model.NewQueueJob(
		MessageName,
		msg,
		append([]model.QueueJobOption{model.QueueJobMaxRetries(2)}, options...)...,
	)
}
