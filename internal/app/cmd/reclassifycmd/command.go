package reclassifycmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen"
)

type Params struct {
	fx.In
	Dao                 lazy.Lazy[*dao.Query]
	ClassifierPublisher lazy.Lazy[publisher.Publisher[message.ClassifyTorrentPayload]]
	Logger              *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	return Result{Command: &cli.Command{
		Name:  "reclassify",
		Usage: "Queue all torrents for reclassification",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "batchSize",
				Value: 100,
			},
		},
		Action: func(ctx *cli.Context) error {
			println("queueing full reclassify...")
			d, err := p.Dao.Get()
			if err != nil {
				return err
			}
			p, err := p.ClassifierPublisher.Get()
			if err != nil {
				return err
			}
			batchSize := ctx.Int("batchSize")
			torrentCount := int64(0)
			if result, err := d.Torrent.WithContext(ctx.Context).Count(); err != nil {
				return err
			} else {
				torrentCount = result
			}
			bar := progressbar.Default(torrentCount, "queuing torrents")
			var torrentResult []*model.Torrent
			if err := d.Torrent.WithContext(ctx.Context).FindInBatches(&torrentResult, batchSize, func(tx gen.Dao, _ int) error {
				infoHashes := make([]protocol.ID, 0, len(torrentResult))
				for _, c := range torrentResult {
					infoHashes = append(infoHashes, c.InfoHash)
				}
				if _, err := p.Publish(ctx.Context, message.ClassifyTorrentPayload{
					InfoHashes: infoHashes,
				}); err != nil {
					return err
				}
				_ = bar.Add(len(torrentResult))
				return nil
			}); err != nil {
				return err
			}
			_ = bar.Finish()
			return nil
		},
	}}, nil
}
