package queue

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"time"
)

const (
	DefaultIdleTxTimeout    = 30000
	DefaultJobCheckInterval = 10 * time.Second
)

// Config configures neoq and its backends
//
// This configuration struct includes options for all backends. As such, some of its options are not applicable to all
// backends. [BackendConcurrency], for example, is only used by the redis backend. Other backends manage concurrency on a
// per-handler basis.
type Config struct {
	JobCheckInterval       time.Duration // the interval of time between checking for new future/retry jobs
	IdleTransactionTimeout int           // the number of milliseconds PgBackend transaction may idle before the connection is killed
	ShutdownTimeout        time.Duration // duration to wait for jobs to finish during shutdown
	SynchronousCommit      bool          // Postgres: Enable synchronous commits (increases durability, decreases performance)
	PGConnectionTimeout    time.Duration // the amount of time to wait for a connection to become available before timing out
}

// ConfigOption is a function that sets optional backend configuration
type ConfigOption func(c *Config)

// NewConfig initiailizes a new Config with defaults
func NewConfig() Config {
	return Config{
		JobCheckInterval: DefaultJobCheckInterval,
	}
}

type Queue interface {
	// Enqueue queues jobs to be executed asynchronously
	Enqueue(ctx context.Context, job model.QueueJob) (jobID string, err error)

	// Start starts processing jobs on the queue specified in the Handler
	Start(ctx context.Context, h handler.Handler) (err error)

	// Shutdown halts job processing and releases resources
	Shutdown(ctx context.Context)
}
