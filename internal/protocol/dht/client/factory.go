package client

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
			logger: p.Logger.Named("dht_client"),
		},
	}
}
