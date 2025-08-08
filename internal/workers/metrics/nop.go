package metrics

import "time"

type nop struct{}

func (nop) IncrAdded(n int) {}

func (nop) IncrDiscarded() {}

func (nop) IncrDeduplicated() {}

func (nop) IncrDequeued(latency time.Duration) {}

func (nop) IncrFlushed(latency time.Duration) {}

func (nop) Reset() {}

var Nop Adapter = nop{}
