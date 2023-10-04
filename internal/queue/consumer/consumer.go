package consumer

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
)

type Consumer interface {
	Pattern() string
	ProcessTask(ctx context.Context, t *asynq.Task) error
}

func New[T interface{}](pattern string, handler Handler[T]) Consumer {
	return &consumer[T]{
		pattern:     pattern,
		unmarshaler: NewJsonUnmarshaler[T](),
		handler:     handler,
	}
}

type Unmarshaler[T interface{}] interface {
	Unmarshal(t *asynq.Task) (T, error)
}

func NewJsonUnmarshaler[T interface{}]() Unmarshaler[T] {
	return &jsonUnmarshaler[T]{}
}

type jsonUnmarshaler[T interface{}] struct{}

func (u jsonUnmarshaler[T]) Unmarshal(t *asynq.Task) (T, error) {
	var payload T
	err := json.Unmarshal(t.Payload(), &payload)
	return payload, err
}

type Handler[T interface{}] interface {
	Handle(ctx context.Context, payload T) error
}

type HandlerFunc[T interface{}] func(ctx context.Context, payload T) error

func (f HandlerFunc[T]) Handle(ctx context.Context, payload T) error {
	return f(ctx, payload)
}

type consumer[T interface{}] struct {
	pattern     string
	unmarshaler Unmarshaler[T]
	handler     Handler[T]
}

func (c consumer[T]) Pattern() string {
	return c.pattern
}

func (c consumer[T]) ProcessTask(ctx context.Context, t *asynq.Task) error {
	payload, err := c.unmarshaler.Unmarshal(t)
	if err != nil {
		return err
	}
	return c.handler.Handle(ctx, payload)
}
