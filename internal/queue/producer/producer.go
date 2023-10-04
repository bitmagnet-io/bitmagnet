package producer

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

type Producer[T interface{}] interface {
	Produce(payload T, options ...asynq.Option) (*asynq.Task, error)
}

func New[T interface{}](typename string, options ...asynq.Option) Producer[T] {
	return &jsonProducer[T]{
		typename: typename,
		options:  options,
	}
}

type jsonProducer[T interface{}] struct {
	typename string
	options  []asynq.Option
}

func (p *jsonProducer[T]) Produce(payload T, options ...asynq.Option) (*asynq.Task, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(p.typename, bytes, append(p.options, options...)...), nil
}
