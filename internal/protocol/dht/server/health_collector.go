package server

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"net/netip"
	"time"
)

type LastResponses struct {
	StartTime    time.Time
	LastSuccess  time.Time
	LastResponse time.Time
}

type healthCollector struct {
	baseServer    Server
	lastResponses *concurrency.AtomicValue[LastResponses]
}

func (c healthCollector) start() error {
	err := c.baseServer.start()
	if err == nil {
		c.lastResponses.Update(func(lr LastResponses) LastResponses {
			lr.StartTime = time.Now()
			return lr
		})
	}
	return err
}

func (c healthCollector) stop() {
	c.lastResponses.Set(LastResponses{})
	c.baseServer.stop()
}

func (c healthCollector) Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (dht.RecvMsg, error) {
	res, err := c.baseServer.Query(ctx, addr, q, args)
	c.lastResponses.Update(func(lr LastResponses) LastResponses {
		lr.LastResponse = time.Now()
		if err == nil {
			lr.LastSuccess = lr.LastResponse
		}
		return lr
	})
	return res, err
}
