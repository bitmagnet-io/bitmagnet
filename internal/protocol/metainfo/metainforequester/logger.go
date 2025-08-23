package metainforequester

import (
	"context"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type requestLogger struct {
	logger    *zap.Logger
	requester Requester
}

func (r *requestLogger) Request(ctx context.Context, infoHash protocol.ID, addr netip.AddrPort) (Response, error) {
	start := time.Now()
	resp, err := r.requester.Request(ctx, infoHash, addr)

	if r.logger.Level() == zapcore.DebugLevel {
		fields := []zap.Field{
			zap.Stringer("info_hash", infoHash),
			zap.Stringer("address", addr),
			zap.Duration("duration", time.Since(start)),
		}

		message := "request"

		if err != nil {
			fields = append(fields, zap.Error(err))
			message += " failed"
		} else {
			fields = append(fields, zap.Stringer("peer_id", resp.PeerID), zap.String("torrent_name", resp.Info.BestName()))
		}

		r.logger.Debug(message, fields...)
	}

	return resp, err
}
