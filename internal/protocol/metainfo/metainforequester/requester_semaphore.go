package metainforequester

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/semaphore"
)

type requesterSemaphore struct {
	Requester
	semaphore semaphore.Semaphore
}

func (r *requesterSemaphore) Request(ctx context.Context, infoHash protocol.ID, addr netip.AddrPort) (Response, error) {
	if err := r.semaphore.Acquire(ctx, 1); err != nil {
		return Response{}, fmt.Errorf("%w: %w: %w", Err, ErrAcquireSemaphore, err)
	}

	defer r.semaphore.Release(1)

	return r.Requester.Request(ctx, infoHash, addr)
}
