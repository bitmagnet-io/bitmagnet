package manager

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
)

type PurgeJobsRequest struct {
	Queues   []string
	Statuses []model.QueueJobStatus
}

type EnqueueReprocessTorrentsBatchRequest struct {
	Purge               bool
	BatchSize           uint
	ChunkSize           uint
	ContentTypes        []model.NullContentType
	Orphans             bool
	ClassifyMode        processor.ClassifyMode
	ClassifierWorkflow  string
	ClassifierFlags     classifier.Flags
	ApisDisabled        bool
	LocalSearchDisabled bool
}

type Manager interface {
	PurgeJobs(context.Context, PurgeJobsRequest) error
	EnqueueReprocessTorrentsBatch(context.Context, EnqueueReprocessTorrentsBatchRequest) error
}
