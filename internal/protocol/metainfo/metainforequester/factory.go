package metainforequester

import (
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net"
	"time"
)

type Params struct {
	fx.In
	Config Config
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Requester Requester
}

func New(p Params) Result {
	return Result{
		Requester: requestLimiter{
			requester: requestLogger{
				requester: requester{
					clientID: protocol.RandomPeerID(),
					timeout:  p.Config.RequestTimeout,
					dialer: &net.Dialer{
						Timeout:   3 * time.Second,
						KeepAlive: -1,
					},
				},
				logger: p.Logger.Named("meta_info_requester"),
			},
			limiter: concurrency.NewKeyedLimiter(rate.Every(time.Second/2), 4, 1000, time.Second*20),
		},
	}
}
