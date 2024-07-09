package dhtcrawler_health_check

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"net"
	"net/netip"
	"time"
)

func NewCheck(
	dhtCrawlerActive *concurrency.AtomicValue[bool],
	lc lazy.Lazy[client.Client],
	bootstrapNodes []string,
) health.Check {
	return health.Check{
		Name: "dht",
		IsActive: func() bool {
			return dhtCrawlerActive.Get()
		},
		Timeout: 60 * 5 * time.Second,
		Check: func(ctx context.Context) error {
			if len(bootstrapNodes) == 0 {
				return errors.New("no bootstrap nodes provided")
			}
			c, cErr := lc.Get()
			if cErr != nil {
				return cErr
			}
			bootstrapNodesRnd := make([]string, len(bootstrapNodes))
			for i := range bootstrapNodes {
				j := rand.Intn(i + 1)
				bootstrapNodesRnd[i], bootstrapNodesRnd[j] = bootstrapNodes[j], bootstrapNodes[i]
			}
			chErrs := make(chan error, len(bootstrapNodes))
			addrs := make(chan netip.AddrPort, len(bootstrapNodes))
			cancelCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			sem := semaphore.NewWeighted(3)
			for _, n := range bootstrapNodesRnd {
				go func() {
					if semErr := sem.Acquire(cancelCtx, 1); semErr != nil {
						chErrs <- semErr
						return
					}
					defer sem.Release(1)
					addr, addrErr := net.ResolveUDPAddr("udp", n)
					if addrErr != nil {
						chErrs <- fmt.Errorf("failed to resolve bootstrap node address: %w", addrErr)
					} else {
						addrs <- addr.AddrPort()
					}
				}()
			}
			success := false
			var errs []error
		outer:
			for {
				select {
				case err := <-chErrs:
					errs = append(errs, err)
					if len(errs) == len(bootstrapNodes) {
						break outer
					}
				case addr := <-addrs:
					_, pingErr := c.Ping(cancelCtx, addr)
					if pingErr != nil {
						chErrs <- pingErr
					} else {
						success = true
						break outer
					}
				}
			}
			if success {
				return nil
			}
			return fmt.Errorf("failed to ping any bootstrap nodes: %w", errors.Join(errs...))
		},
	}
}
