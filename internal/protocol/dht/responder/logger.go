package responder

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"go.uber.org/zap"
)

type responderLogger struct {
	responder Responder
	logger    *zap.Logger
}

func (r responderLogger) Respond(ctx context.Context, msg dht.RecvMsg) (dht.Return, error) {
	start := time.Now()
	ret, err := r.responder.Respond(ctx, msg)

	fields := []zap.Field{
		zap.Stringer("from", msg.From),
		zap.Duration("duration", time.Since(start)),
	}

	message := msg.Msg.Q

	if err == nil {
		switch msg.Msg.Q {
		case dht.QFindNode:
			fields = append(fields,
				zap.Stringer("target", msg.Msg.A.Target),
				zap.Int("nodes", len(ret.Nodes)),
			)
		case dht.QGetPeers:
			fields = append(fields,
				zap.Int("values", len(ret.Values)),
				zap.Int("nodes", len(ret.Nodes)),
				zap.Stringp("token", ret.Token),
			)
		case dht.QAnnouncePeer:
			fields = append(fields,
				zap.Stringer("info_hash", msg.Msg.A.InfoHash),
				zap.Uint16("port", msg.AnnouncePort()),
				zap.String("token", msg.Msg.A.Token))
		case dht.QSampleInfohashes:
			fields = append(fields,
				zap.Int("samples", len(*ret.Samples)),
				zap.Int("nodes", len(ret.Nodes)),
				zap.Int64("num", *ret.Num),
				zap.Int64("interval", *ret.Interval))
		}
	} else {
		message += " error"

		fields = append(fields, zap.Error(err))
	}

	r.logger.Debug(message, fields...)

	return ret, err
}
