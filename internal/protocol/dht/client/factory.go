package client

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Params struct {
	fx.In
	NodeID protocol.ID `name:"dht_node_id"`
	Server server.Server
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Client Client
}

func New(p Params) Result {
	return Result{
		Client: clientLogger{
			client: serverAdapter{
				nodeID: p.NodeID,
				server: p.Server,
			},
			// we make way to many queries to usefully log everything, but having a sample is helpful:
			logger: p.Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Minute, 10, 0)
			})).Named("dht_client"),
		},
	}
}
