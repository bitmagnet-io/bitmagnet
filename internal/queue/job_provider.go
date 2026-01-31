package queue

import "github.com/bitmagnet-io/bitmagnet/internal/model"

type JobProvider[Msg any] func(msg Msg, options ...model.QueueJobOption) (model.QueueJob, error)
