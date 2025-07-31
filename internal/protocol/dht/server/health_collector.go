package server

import (
	"context"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
)

type LastResponses struct {
	StartTime    time.Time
	LastSuccess  time.Time
	LastResponse time.Time
}

type healthCollector struct {
	serverRunner
	lastResponses *concurrency.AtomicValue[LastResponses]
}

func (c healthCollector) Query(
	ctx context.Context,
	addr netip.AddrPort,
	q string,
	args dht.MsgArgs,
) (dht.RecvMsg, error) {
	res, err := c.serverRunner.Query(ctx, addr, q, args)
	c.lastResponses.Update(func(lr LastResponses) LastResponses {
		lr.LastResponse = time.Now()
		if err == nil {
			lr.LastSuccess = lr.LastResponse
		}

		return lr
	})

	return res, err
}
