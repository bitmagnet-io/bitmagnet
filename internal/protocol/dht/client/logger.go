package client

import (
	"context"
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
	l.logger.Debugw(dht.QPing, "addr", addr, "duration", time.Since(start), "error", err)
	return res, err
}

func (l clientLogger) FindNode(ctx context.Context, addr netip.AddrPort, target dht.ID) (FindNodeResult, error) {
	start := time.Now()
	res, err := l.client.FindNode(ctx, addr, target)
	l.logger.Debugw(dht.QFindNode, "addr", addr, "target", target, "duration", time.Since(start), "nodes", len(res.Nodes), "error", err)
	return res, err
}

func (l clientLogger) GetPeers(ctx context.Context, addr netip.AddrPort, infoHash dht.ID) (GetPeersResult, error) {
	start := time.Now()
	res, err := l.client.GetPeers(ctx, addr, infoHash)
	l.logger.Debugw(dht.QGetPeers, "addr", addr, "infoHash", infoHash, "duration", time.Since(start), "values", len(res.Values), "nodes", len(res.Nodes), "error", err)
	return res, err
}

func (l clientLogger) GetPeersScrape(ctx context.Context, addr netip.AddrPort, infoHash dht.ID) (GetPeersScrapeResult, error) {
	start := time.Now()
	res, err := l.client.GetPeersScrape(ctx, addr, infoHash)
	l.logger.Debugw(dht.QGetPeers, "addr", addr, "infoHash", infoHash, "duration", time.Since(start), "error", err)
	return res, err
}

func (l clientLogger) SampleInfoHashes(ctx context.Context, addr netip.AddrPort, target dht.ID) (SampleInfoHashesResult, error) {
	start := time.Now()
	res, err := l.client.SampleInfoHashes(ctx, addr, target)
	l.logger.Debugw(dht.QSampleInfohashes, "addr", addr, "target", target, "duration", time.Since(start), "samples", len(res.Samples), "num", res.Num, "interval", res.Interval, "nodes", len(res.Nodes), "error", err)
	return res, err
}
