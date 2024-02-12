package publisher

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
)

type Publisher[T interface{}] interface {
	Publish(ctx context.Context, payload T) (model.QueueJob, error)
}

func New[T interface{}](client queue.Queue, queueName string) Publisher[T] {
	return &publisher[T]{
		client:    client,
		queueName: queueName,
		//producer: producer,
	}
}

type publisher[T interface{}] struct {
	client    queue.Queue
	queueName string
	//producer producer.Producer[T]
}

func (p publisher[T]) Publish(ctx context.Context, payload T) (model.QueueJob, error) {
	//task, err := p.producer.Produce(payload)
	//if err != nil {
	//	return nil, err
	//}
	//// todo: Add ID!
	job, err := model.NewQueueJob(p.queueName, payload)
	if err != nil {
		return model.QueueJob{}, err
	}
	_, err = p.client.Enqueue(ctx, job)
	if err != nil {
		return model.QueueJob{}, err
	}
	//task.ID = id
	return job, nil
}
