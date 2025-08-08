package metrics

import (
	"time"
)

// Snapshot represents a point-in-time capture of all metrics in the registry.
type Snapshot struct {
	Totals
	Averages Averages
	InitTime time.Time
	Duration time.Duration
}

// Snapshot returns a Snapshot of the current state of all metrics.
func (r *Registry) Snapshot() Snapshot {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return Snapshot{
		Totals:   r.sink.totalsUnlocked(),
		Averages: r.sink.averagesUnlocked(),
		InitTime: r.startTime,
		Duration: time.Since(r.startTime),
	}
}
