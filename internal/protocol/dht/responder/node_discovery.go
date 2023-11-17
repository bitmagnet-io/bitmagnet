package responder

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

// responderNodeDiscovery attempts to add nodes from incoming requests to the discovered nodes channel.
type responderNodeDiscovery struct {
	responder       Responder
	discoveredNodes chan<- ktable.Node
}

func (r responderNodeDiscovery) Respond(ctx context.Context, msg dht.RecvMsg) (dht.Return, error) {
	ret, err := r.responder.Respond(ctx, msg)
	if err == nil {
		go func() {
			// wait for up to a second
			cancelCtx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			select {
			case <-cancelCtx.Done():
			case r.discoveredNodes <- ktable.NewNode(msg.Msg.A.ID, msg.From):
			}
		}()
	}
	return ret, err
}
