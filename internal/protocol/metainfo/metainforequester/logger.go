package metainforequester

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"go.uber.org/zap"
	"net/netip"
	"time"
)

type requestLogger struct {
	logger    *zap.SugaredLogger
	requester Requester
}

func (r requestLogger) Request(ctx context.Context, infoHash protocol.ID, addr netip.AddrPort) (Response, error) {
	start := time.Now()
	resp, err := r.requester.Request(ctx, infoHash, addr)
	keyValues := []interface{}{
		"infoHash", infoHash,
		"addr", addr,
		"duration", time.Since(start),
	}
	message := "request"
	if err != nil {
		keyValues = append(keyValues, "error", err)
		message += " failed"
	} else {
		keyValues = append(keyValues,
			"peerId", resp.PeerID,
			"torrentName", resp.Info.BestName(),
		)
	}
	r.logger.Debugw(message, keyValues...)
	return resp, err
}
