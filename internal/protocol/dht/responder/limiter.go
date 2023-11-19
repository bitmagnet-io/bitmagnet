package responder

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"golang.org/x/time/rate"
	"net/netip"
	"time"
)

// responderLimiter applies both overall and per-IP rate limiting
type responderLimiter struct {
	responder Responder
	limiter   Limiter
}

func (r responderLimiter) Respond(ctx context.Context, msg dht.RecvMsg) (ret dht.Return, err error) {
	if !r.limiter.Allow(msg.From.Addr()) {
		return dht.Return{}, ErrTooManyRequests
	}
	return r.responder.Respond(ctx, msg)
}

type Limiter interface {
	Allow(addr netip.Addr) bool
}

type limiter struct {
	keyedLimiter concurrency.KeyedLimiter
	limiter      *rate.Limiter
}

func NewLimiter(
	overallRate rate.Limit,
	overallBurst int,
	perIpRate rate.Limit,
	perIpBurst int,
	perIpSize int,
	perIpTtl time.Duration,
) Limiter {
	return &limiter{
		limiter:      rate.NewLimiter(overallRate, overallBurst),
		keyedLimiter: concurrency.NewKeyedLimiter(perIpRate, perIpBurst, perIpSize, perIpTtl),
	}
}

func (l *limiter) Allow(addr netip.Addr) bool {
	return l.keyedLimiter.Allow(addr.String()) && l.limiter.Allow()
}
