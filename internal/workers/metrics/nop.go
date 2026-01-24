package metrics

import "time"

type nop struct{}

func (nop) IncrAdded(int) {}

func (nop) IncrDiscarded() {}

func (nop) IncrDeduplicated() {}

func (nop) IncrDequeued(time.Duration) {}

func (nop) IncrFlushed(time.Duration) {}

func (nop) Reset() {}

var Nop Adapter = nop{}
