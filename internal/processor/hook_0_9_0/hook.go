package hook_0_9_0

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	"go.uber.org/fx"
	"gorm.io/gorm/clause"
	"time"
)

type Params struct {
	fx.In
	Dao lazy.Lazy[*dao.Query]
}

type Result struct {
	fx.Out
	Dao lazy.Lazy[*dao.Query]
}

func NewDecorator(p Params) Result {
	return Result{
		Dao: lazy.New(func() (*dao.Query, error) {
			d, dErr := p.Dao.Get()
			if dErr != nil {
				return nil, dErr
			}
			if execErr := execHook(d); execErr != nil {
				return nil, execErr
			}
			return d, nil
		}),
	}
}

const keyName = "0.9.0_reprocess_enqueued"

func execHook(d *dao.Query) error {
	ctx := context.Background()
	return d.Transaction(func(tx *dao.Query) error {
		if result, err := tx.KeyValue.
			WithContext(ctx).
			Where(d.KeyValue.Key.Eq(keyName)).
			Limit(1).
			Find(); err != nil {
			return err
		} else if len(result) > 0 {
			return nil
		}
		if err := tx.KeyValue.
			WithContext(ctx).
			Create(&model.KeyValue{
				Key:   keyName,
				Value: "1",
			}); err != nil {
			return err
		}
		job, jobErr := batch.NewQueueJob(batch.MessageParams{
			UpdatedBefore: time.Now(),
			ClassifierFlags: classifier.Flags{
				"apis_enabled": false,
			},
			ChunkSize: 10_000,
			BatchSize: 100,
		})
		if jobErr != nil {
			return jobErr
		}
		return tx.QueueJob.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&job)
	})
}
