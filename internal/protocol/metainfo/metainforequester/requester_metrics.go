package metainforequester

import (
	"context"
	"net/netip"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

func newRequesterMetrics(
	requester Requester,
	component *metrics.Component,
) Requester {
	return &requesterMetrics{
		Requester:        requester,
		sampler:          component.MustNewSampler("requests"),
		concurrencyGauge: component.MustNewGauge("concurrency"),
	}
}

type requesterMetrics struct {
	Requester
	mtx              sync.Mutex
	concurrency      int
	sampler          *metrics.Sampler
	concurrencyGauge *metrics.Gauge
}

var (
	labelStatus       = metrics.Label("status")
	labelValueSuccess = labelStatus.Value("success")
	labelValueError   = labelStatus.Value("error")
)

func (r *requesterMetrics) Request(ctx context.Context, infoHash protocol.ID, addr netip.AddrPort) (Response, error) {
	r.mtx.Lock()
	r.concurrency++
	r.concurrencyGauge.Set(r.concurrency)
	r.mtx.Unlock()

	startTime := time.Now()

	res, err := r.Requester.Request(ctx, infoHash, addr)

	duration := time.Since(startTime)

	r.mtx.Lock()
	r.concurrency--
	r.concurrencyGauge.Set(r.concurrency)
	r.mtx.Unlock()

	label := labelValueSuccess
	if err != nil {
		label = labelValueError
	}

	r.sampler.Add(float32(duration.Seconds()), label)

	return res, err
}
