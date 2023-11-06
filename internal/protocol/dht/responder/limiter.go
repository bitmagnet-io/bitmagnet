package responder

import (
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"golang.org/x/time/rate"
	"net/netip"
	"time"
)

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
