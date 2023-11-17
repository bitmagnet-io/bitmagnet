package client

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"go.uber.org/zap"
	"net/netip"
	"time"
)

type clientLogger struct {
	client Client
	logger *zap.SugaredLogger
}

func (l clientLogger) Ready() <-chan struct{} {
	return l.client.Ready()
}

func (l clientLogger) Ping(ctx context.Context, addr netip.AddrPort) (PingResult, error) {
	start := time.Now()
	res, err := l.client.Ping(ctx, addr)
	l.log(dht.QPing, addr, start, err)
	return res, err
}

func (l clientLogger) FindNode(ctx context.Context, addr netip.AddrPort, target dht.ID) (FindNodeResult, error) {
	start := time.Now()
	res, err := l.client.FindNode(ctx, addr, target)
	l.log(dht.QFindNode, addr, start, err)
	return res, err
}

func (l clientLogger) GetPeers(ctx context.Context, addr netip.AddrPort, infoHash dht.ID) (GetPeersResult, error) {
	start := time.Now()
	res, err := l.client.GetPeers(ctx, addr, infoHash)
	l.log(dht.QGetPeers, addr, start, err)
	return res, err
}

func (l clientLogger) GetPeersScrape(ctx context.Context, addr netip.AddrPort, infoHash dht.ID) (GetPeersScrapeResult, error) {
	start := time.Now()
	res, err := l.client.GetPeersScrape(ctx, addr, infoHash)
	l.log(dht.QGetPeers+":scrape", addr, start, err)
	return res, err
}

func (l clientLogger) SampleInfoHashes(ctx context.Context, addr netip.AddrPort, target dht.ID) (SampleInfoHashesResult, error) {
	start := time.Now()
	res, err := l.client.SampleInfoHashes(ctx, addr, target)
	l.log(dht.QSampleInfohashes, addr, start, err)
	return res, err
}

func (l clientLogger) log(query string, addr netip.AddrPort, start time.Time, err error) {
	if err == nil {
		l.logger.Debugw(query, "addr", addr, "duration", time.Since(start))
	} else {
		l.logger.Debugw(fmt.Sprintf("%s failed", query), "addr", addr, "duration", time.Since(start), "error", err)
	}
}
