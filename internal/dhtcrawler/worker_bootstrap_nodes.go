package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"net"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/periodic"
)

func newBootstrapNodesWorker(
	interval time.Duration,
	bootstrapNodes []string,
	discoveredNodesAdder channel.Adder[ktable.Node],
	logger logging.Logger,
) runner.Provider {
	return periodic.New(
		interval,
		func(ctx context.Context) error {
			for _, strAddr := range bootstrapNodes {
				addr, err := net.ResolveUDPAddr("udp", strAddr)
				if err != nil {
					logger.Warnf("failed to resolve bootstrap node address: %s", err)
					continue
				}

				err = discoveredNodesAdder.Add(ctx, ktable.NewNode(ktable.ID{}, addr.AddrPort()))
				if err != nil {
					return err
				}
			}

			return nil
		},
		periodic.WithInitialInterval(0),
		periodic.WithQuickShutdown(),
	)
}
