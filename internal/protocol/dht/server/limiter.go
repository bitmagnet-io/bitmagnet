package server

import (
	"context"
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
)

type queryLimiter struct {
	serverRunner
	queryLimiter concurrency.KeyedLimiter
}

func (s queryLimiter) Query(
	ctx context.Context,
	addr netip.AddrPort,
	q string,
	args dht.MsgArgs,
) (r dht.RecvMsg, err error) {
	if limitErr := s.queryLimiter.Wait(ctx, addr.Addr().String()); limitErr != nil {
		return r, limitErr
	}

	return s.serverRunner.Query(ctx, addr, q, args)
}
