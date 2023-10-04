package publisher

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/producer"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Client   *asynq.Client
	Producer producer.Producer[message.ClassifyTorrentPayload]
}

type Result struct {
	fx.Out
	Publisher publisher.Publisher[message.ClassifyTorrentPayload]
}

func New(p Params) Result {
	return Result{
		Publisher: publisher.New[message.ClassifyTorrentPayload](p.Client, p.Producer),
	}
}
