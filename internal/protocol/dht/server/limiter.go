package server

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"net/netip"
)

type queryLimiter struct {
	server       Server
	queryLimiter concurrency.KeyedLimiter
}

func (s queryLimiter) Ready() <-chan struct{} {
	return s.server.Ready()
}

func (s queryLimiter) Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (r dht.RecvMsg, err error) {
	if limitErr := s.queryLimiter.Wait(ctx, addr.Addr().String()); limitErr != nil {
		return r, limitErr
	}
	return s.server.Query(ctx, addr, q, args)
}
