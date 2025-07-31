package metainforequester

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"net/netip"
)

type requesterSemaphore struct {
	requester Requester
	semaphore chan struct{}
}

func (r *requesterSemaphore) Request(ctx context.Context, infoHash protocol.ID, node netip.AddrPort) (Response, error) {
	select {
	case <-ctx.Done():
		return Response{}, ctx.Err()
	case r.semaphore <- struct{}{}:
	}

	defer func() { <-r.semaphore }()

	return r.requester.Request(ctx, infoHash, node)
}
