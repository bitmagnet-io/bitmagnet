package httpserver

import (
	"context"
	"net"
	"net/http"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func New(
	handler http.Handler,
	localAddress LocalAddress,
) runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		srv := &http.Server{
			Addr:    string(localAddress),
			Handler: handler,
		}

		var lc net.ListenConfig

		ln, listenErr := lc.Listen(ctx, "tcp", srv.Addr)
		if listenErr != nil {
			return runner.NopShutdowner, listenErr
		}

		shutdown := make(chan struct{})

		go func() {
			err := srv.Serve(ln)

			select {
			case <-shutdown:
			default:
				cancel(err)
			}
		}()

		return func(ctx context.Context) error {
			close(shutdown)
			return srv.Shutdown(ctx)
		}, nil
	}
}
