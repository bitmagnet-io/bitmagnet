package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/server"
	"github.com/hellofresh/health-go/v5"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"net"
	"sync"
	"time"
)

type Params struct {
	fx.In
	DhtServer server.Server
}

type Result struct {
	fx.Out
	Option health.Option `group:"healthcheck_options"`
}

func New(p Params) Result {
	return Result{
		Option: health.WithChecks(health.Config{
			Name:    "dht",
			Timeout: time.Second * 15,
			Check:   newHealthCheck(p.DhtServer),
		}),
	}
}

// newHealthCheck returns a DHT server health check function.
// This is quite a complex health check:
// We expect the bootstrap nodes to be flaky, so we only want a successful response from one of them to consider ourselves healthy.
// We'll try pinging 3 nodes in parallel, stopping when we get a successful response.
// Each node will have 2 seconds to respond before we move on to try the next node.
func newHealthCheck(s server.Server) health.CheckFunc {
	// make a map to get random ordering for free
	bootstrapNodes := make(map[string]struct{})
	for _, n := range dht.DefaultGlobalBootstrapHostPorts {
		bootstrapNodes[n] = struct{}{}
	}
	return func(ctx context.Context) error {
		checkCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		success := make(chan struct{})
		done := make(chan struct{})
		sem := semaphore.NewWeighted(3)
		errsMutex := &sync.Mutex{}
		var errs []error
		addErr := func(err error) {
			errsMutex.Lock()
			errs = append(errs, err)
			errsMutex.Unlock()
		}
		wg := sync.WaitGroup{}
		wg.Add(len(bootstrapNodes))
		go func() {
			wg.Wait()
			close(done)
		}()
		for node := range bootstrapNodes {
			go func(node string) {
				defer wg.Done()
				if semErr := sem.Acquire(checkCtx, 1); semErr != nil {
					addErr(fmt.Errorf("failed to acquire semaphore: %w", semErr))
					return
				}
				defer sem.Release(1)
				if ctxErr := checkCtx.Err(); ctxErr != nil {
					addErr(ctxErr)
					return
				}
				addr, resolveErr := net.ResolveUDPAddr("udp", node)
				if resolveErr != nil {
					addErr(fmt.Errorf("failed to resolve UDP address for '%s': %w", node, resolveErr))
					return
				}
				// we'll only give each node 2 seconds to respond to avoid delaying the health check
				thisCtx, thisCancel := context.WithTimeout(checkCtx, 2*time.Second)
				_, pingErr := s.QueryUrgent(thisCtx, krpc.NodeAddr{
					IP:   addr.IP.To4(),
					Port: addr.Port,
				}, "ping", krpc.MsgArgs{})
				thisCancel()
				if pingErr == nil {
					close(success)
					return
				}
				addErr(fmt.Errorf("failed to ping bootstrap node '%s': %w", node, pingErr))
			}(node)
		}
		select {
		case <-success:
			return nil
		case <-done:
			err := errors.Join(errs...)
			if err == nil {
				err = errors.New("unknown error")
			}
			return fmt.Errorf("failed to ping any bootstrap node: %w", err)
		}
	}
}
