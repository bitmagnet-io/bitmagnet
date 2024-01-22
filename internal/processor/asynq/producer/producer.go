package producer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/producer"
	"github.com/hibiken/asynq"
	"time"
)

func New() producer.Producer[processor.MessageParams] {
	return producer.New[processor.MessageParams](
		processor.MessageName,
		asynq.Queue(processor.MessageName),
		asynq.MaxRetry(1),
		// high retention here allows for a large queue to be fully worked down
		asynq.Retention(time.Hour*24*7),
		asynq.Unique(time.Hour*24),
	)
}
