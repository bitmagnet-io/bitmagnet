package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func New(proc processor.Processor) lazy.Lazy[handler.Handler] {
	return lazy.New(func() (handler.Handler, error) {
		return handler.New(
			processor.MessageName,
			func(job model.QueueJob) runner.Runner {
				return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
					msg := &processor.MessageParams{}
					if err := json.Unmarshal([]byte(job.Payload), msg); err != nil {
						return runner.NopShutdowner, err
					}

					return proc.NewJob(*msg)(ctx, cancel)
				}
			},
			handler.JobTimeout(time.Second*60*10),
			handler.Concurrency(1),
		), nil
	})
}
