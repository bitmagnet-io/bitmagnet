package metainforequester

import (
	"context"
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type requestLimiter struct {
	requester Requester
	limiter   concurrency.KeyedLimiter
}

func (r requestLimiter) Request(ctx context.Context, infoHash protocol.ID, node netip.AddrPort) (Response, error) {
	if limitErr := r.limiter.Wait(ctx, node.Addr().String()); limitErr != nil {
		return Response{}, limitErr
	}
	return r.requester.Request(ctx, infoHash, node)
}
