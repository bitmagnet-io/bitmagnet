package dht_crawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/batch"
)

func newEnqueueProcessTorrentWorker(
	queueJobProvider queue.JobProvider[processor.MessageParams],
	persisterAdder persister.Adder,
) batch.Worker[protocol.ID] {
	return batch.NewWorker(
		batch.WithValuesAsKeys[protocol.ID](),
		batch.WithMaxSize[protocol.ID, protocol.ID](100),
		batch.WithMaxWait[protocol.ID, protocol.ID](time.Minute),
		batch.WithFlusher[protocol.ID](func(ctx context.Context, infoHashes []protocol.ID) error {
			job, err := queueJobProvider(processor.MessageParams{
				InfoHashes: infoHashes,
			},
				// delay the classifier by a minute to allow time for the S/L scrape:
				model.QueueJobDelayBy(time.Minute),
			)
			if err != nil {
				return err
			}

			return persisterAdder.Add(ctx, persister.InputQueueJobs(job))
		}),
	)
}
