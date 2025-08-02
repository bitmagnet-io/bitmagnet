package client

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(
	nodeID protocol.ID,
	server server.Server,
	logger *zap.Logger,
) Client {
	return clientLogger{
		client: serverAdapter{
			nodeID: nodeID,
			server: server,
		},
		// we make way to many queries to usefully log everything, but having a sample is
		// helpful:
		logger: logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(core, time.Minute, 10, 0)
		})).Named("dht_client"),
	}
}
