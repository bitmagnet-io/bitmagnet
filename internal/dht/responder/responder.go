package responder

import (
	"context"
	"errors"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config dht.Config
}

type Result struct {
	fx.Out
	Responder Responder
}

func New(p Params) Result {
	return Result{
		Responder: responder{
			peerID: p.Config.PeerID,
		},
	}
}

type Responder interface {
	Respond(context.Context, krpc.Msg) (krpc.Return, error)
}

type responder struct {
	peerID [20]byte
}

var ErrUnsupportedQuery = errors.New("unsupported query")

func (r responder) Respond(_ context.Context, msg krpc.Msg) (ret krpc.Return, _ error) {
	ret.ID = r.peerID
	switch msg.Q {
	case "ping":
		return ret, nil
	}
	return ret, ErrUnsupportedQuery
}
