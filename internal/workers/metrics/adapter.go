package metrics

import "time"

type Adapter interface {
	IncrAdded(n int)
	IncrDiscarded()
	IncrDeduplicated()
	IncrDequeued(latency time.Duration)
	IncrFlushed(latency time.Duration)
	Reset()
}
