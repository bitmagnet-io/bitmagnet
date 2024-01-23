// This file contains a hook for the 0.5.0 release that enqueues a full reprocess of all torrents.
// This is a one-time job. The decorator will be removed in a later version.

package decorator

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	Dao       lazy.Lazy[*dao.Query]
	Publisher lazy.Lazy[publisher.Publisher[processor.MessageParams]]
	Logger    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Decorator worker.Decorator `group:"worker_decorators"`
}

const checkKeyValue = "0.5.0_reprocess_enqueued"

func New(p Params) (Result, error) {
	logger := p.Logger.Named("0.5.0_upgrade")
	return Result{
		Decorator: worker.Decorator{
			Key: "queue_server",
			Decorate: func(hook fx.Hook) fx.Hook {
				return fx.Hook{
					OnStart: func(ctx context.Context) error {
						d, daoErr := p.Dao.Get()
						if daoErr != nil {
							return daoErr
						}
						enqueued := false
						if kvErr := d.Transaction(func(tx *dao.Query) error {
							_, err := tx.KeyValue.WithContext(ctx).Where(tx.KeyValue.Key.Eq(checkKeyValue)).First()
							if err != nil {
								if !errors.Is(err, gorm.ErrRecordNotFound) {
									return err
								}
								return tx.KeyValue.WithContext(ctx).Create(&model.KeyValue{
									Key:   checkKeyValue,
									Value: "1",
								})
							}
							enqueued = true
							return nil
						}); kvErr != nil {
							return kvErr
						}
						if enqueued {
							logger.Debug("reprocess already enqueued, skipping")
						} else {
							logger.Info("enqueuing reprocess")
							p, pubErr := p.Publisher.Get()
							if pubErr != nil {
								return pubErr
							}
							torrentCount := int64(0)
							if result, err := d.Torrent.WithContext(ctx).Count(); err != nil {
								return err
							} else {
								torrentCount = result
							}
							if torrentCount > 0 {
								batchSize := 100
								bar := progressbar.Default(torrentCount, "queuing torrents")
								var torrentResult []*model.Torrent
								if err := d.Torrent.WithContext(ctx).FindInBatches(&torrentResult, batchSize, func(tx gen.Dao, _ int) error {
									infoHashes := make([]protocol.ID, 0, len(torrentResult))
									for _, c := range torrentResult {
										infoHashes = append(infoHashes, c.InfoHash)
									}
									if _, err := p.Publish(ctx, processor.MessageParams{
										InfoHashes: infoHashes,
									}); err != nil {
										return err
									}
									_ = bar.Add(len(torrentResult))
									return nil
								}); err != nil {
									return err
								}
							}
							logger.Info("enqueued reprocess")
						}
						return hook.OnStart(ctx)
					},
					OnStop: hook.OnStop,
				}
			},
		},
	}, nil
}
