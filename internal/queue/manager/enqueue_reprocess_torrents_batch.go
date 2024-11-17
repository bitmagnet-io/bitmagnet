package manager

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	"gorm.io/gorm"
	"time"
)

func (m manager) EnqueueReprocessTorrentsBatch(ctx context.Context, req EnqueueReprocessTorrentsBatchRequest) error {
	flags := req.ClassifierFlags
	if flags == nil {
		flags = make(classifier.Flags)
	}
	if req.ApisDisabled {
		flags["apis_enabled"] = false
	}
	if req.LocalSearchDisabled {
		flags["local_search_enabled"] = false
	}
	job, err := batch.NewQueueJob(batch.MessageParams{
		ClassifyMode:    req.ClassifyMode,
		ClassifierFlags: flags,
		ChunkSize:       req.ChunkSize,
		BatchSize:       req.BatchSize,
		ContentTypes:    req.ContentTypes,
		Orphans:         req.Orphans,
		UpdatedBefore:   time.Now(),
	})
	if err != nil {
		return err
	}
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if req.Purge {
			if _, err := tx.WithContext(ctx).Raw("TRUNCATE TABLE queue_jobs;").Rows(); err != nil {
				return fmt.Errorf("error purging queue: %w", err)
			}
		}
		return tx.WithContext(ctx).Create(&job).Error
	})
}
