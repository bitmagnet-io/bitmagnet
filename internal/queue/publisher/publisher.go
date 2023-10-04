package publisher

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/producer"
	"github.com/hibiken/asynq"
)

type Publisher[T interface{}] interface {
	Publish(ctx context.Context, payload T, options ...asynq.Option) (*asynq.TaskInfo, error)
}

func New[T interface{}](client *asynq.Client, producer producer.Producer[T]) Publisher[T] {
	return &publisher[T]{
		client:   client,
		producer: producer,
	}
}

type publisher[T interface{}] struct {
	client   *asynq.Client
	producer producer.Producer[T]
}

func (p publisher[T]) Publish(ctx context.Context, payload T, options ...asynq.Option) (*asynq.TaskInfo, error) {
	task, err := p.producer.Produce(payload, options...)
	if err != nil {
		return nil, err
	}
	return p.client.EnqueueContext(ctx, task, options...)
}
