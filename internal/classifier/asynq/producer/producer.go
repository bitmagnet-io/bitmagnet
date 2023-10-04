package producer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/producer"
	"github.com/hibiken/asynq"
	"time"
)

func New() producer.Producer[message.ClassifyTorrentPayload] {
	return producer.New[message.ClassifyTorrentPayload](
		message.ClassifyTorrentTypename,
		asynq.Queue(message.ClassifyTorrentTypename),
		asynq.MaxRetry(1),
		// high retention here allows for a large queue to be fully worked down
		asynq.Retention(time.Hour*24*7),
		asynq.Unique(time.Hour*24),
	)
}
