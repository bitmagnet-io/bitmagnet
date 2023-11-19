package responder

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"go.uber.org/zap"
	"time"
)

type responderLogger struct {
	responder Responder
	logger    *zap.SugaredLogger
}

func (r responderLogger) Respond(ctx context.Context, msg dht.RecvMsg) (dht.Return, error) {
	start := time.Now()
	ret, err := r.responder.Respond(ctx, msg)
	var logData []interface{}
	log := func(k string, v interface{}) {
		logData = append(logData, k, v)
	}
	message := msg.Msg.Q
	log("from", msg.From)
	log("duration", time.Since(start))
	if err == nil {
		switch msg.Msg.Q {
		case dht.QFindNode:
			log("target", msg.Msg.A.Target)
			log("nodes", len(ret.Nodes))
		case dht.QGetPeers:
			log("values", len(ret.Values))
			log("nodes", len(ret.Nodes))
			log("token", ret.Token)
		case dht.QAnnouncePeer:
			log("infoHash", msg.Msg.A.InfoHash)
			log("port", msg.AnnouncePort())
			log("token", msg.Msg.A.Token)
		case dht.QSampleInfohashes:
			log("samples", len(*ret.Samples))
			log("nodes", len(ret.Nodes))
			log("num", *ret.Num)
			log("interval", *ret.Interval)
		}
	} else {
		message += " error"
		log("error", err)
	}
	r.logger.Debugw(message, logData...)
	return ret, err
}
