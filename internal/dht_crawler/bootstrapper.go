package dht_crawler

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
)

func newBootstrapper(
	bootstrapNodes []string,
	discoveredNodesAdder channel.Adder[ktable.Node],
) func(context.Context) error {
	return func(ctx context.Context) error {
		var (
			wg      sync.WaitGroup
			mtx     sync.Mutex
			success bool
			errs    []error
		)

		addError := func(err error) {
			mtx.Lock()
			defer mtx.Unlock()

			errs = append(errs, err)
		}

		for _, strAddr := range bootstrapNodes {
			wg.Add(1)

			go func(strAddr string) {
				defer wg.Done()

				addr, err := net.ResolveUDPAddr("udp", strAddr)
				if err != nil {
					addError(fmt.Errorf("failed to resolve bootstrap node: %s: %w", strAddr, err))
					return
				}

				err = discoveredNodesAdder.Add(ctx, ktable.NewNode(ktable.ID{}, addr.AddrPort()))
				if err != nil {
					addError(fmt.Errorf("failed to add bootstrap node: %s: %w", strAddr, err))
					return
				}

				success = true
			}(strAddr)
		}

		wg.Wait()

		if !success {
			return fmt.Errorf("bootstrap failed: %w", errors.Join(errs...))
		}

		return nil
	}
}
