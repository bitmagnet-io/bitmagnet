package producer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Producer[T interface{}] interface {
	Produce(
		payload T,
		//options ...asynq.Option,
	) (*model.QueueJob, error)
}

//
//func New[T interface{}](typename string) Producer[T] {
//	return &jsonProducer[T]{
//		typename: typename,
//		//options:  options,
//	}
//}
//
//type jsonProducer[T interface{}] struct {
//	typename string
//	//options  []asynq.Option
//}
//
//func (p *jsonProducer[T]) Produce(payload T) (*model.QueueJob, error) {
//	bytes, err := json.Marshal(payload)
//	if err != nil {
//		return nil, err
//	}
//	j := &jobs.Job{
//		Queue:      p.typename,
//		Payload:    string(bytes),
//		MaxRetries: 2,
//	}
//	if err := jobs.FingerprintJob(j); err != nil {
//		return nil, err
//	}
//	return j, nil
//}
