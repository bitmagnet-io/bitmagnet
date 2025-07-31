package responder

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
)

// responderNodeDiscovery attempts to add nodes from incoming requests to the discovered nodes channel.
type responderNodeDiscovery struct {
	responder       Responder
	discoveredNodes channel.Worker[ktable.Node]
}

func (r responderNodeDiscovery) Respond(ctx context.Context, msg dht.RecvMsg) (dht.Return, error) {
	ret, err := r.responder.Respond(ctx, msg)
	if err == nil {
		go func() {
			// wait for up to a second
			cancelCtx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			_ = r.discoveredNodes.Add(cancelCtx, ktable.NewNode(msg.Msg.A.ID, msg.From))
		}()
	}

	return ret, err
}
