package queue

import (
	"context"
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"gorm.io/gen"
)

func New(
	daoProvider database.DaoProvider,
	processJobProvider queue.JobProvider[processor.MessageParams],
	batchJobProvider queue.JobProvider[batch.MessageParams],
) handler.Func {
	return func(job model.QueueJob) runner.Runner {
		return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
			msg := &batch.MessageParams{}
			if err := json.Unmarshal([]byte(job.Payload), msg); err != nil {
				return runner.NopShutdowner, err
			}

			dao, err := daoProvider.Dao()
			if err != nil {
				return runner.NopShutdowner, err
			}

			var scopes []func(gen.Dao) gen.Dao

			if len(msg.ContentTypes) > 0 {
				var contentTypes []string

				var unknownContentType bool

				for _, ct := range msg.ContentTypes {
					if !ct.Valid {
						unknownContentType = true
					} else {
						contentTypes = append(contentTypes, ct.ContentType.String())
					}
				}

				scopes = append(scopes, func(tx gen.Dao) gen.Dao {
					sq := dao.TorrentContent.Where(
						dao.TorrentContent.InfoHash.EqCol(dao.Torrent.InfoHash),
					).Where(dao.TorrentContent.ContentType.In(contentTypes...))
					if unknownContentType {
						sq = sq.Or(dao.TorrentContent.ContentType.IsNull())
					}

					return tx.Where(gen.Exists(sq))
				})
			}

			if msg.Orphans {
				scopes = append(scopes, func(tx gen.Dao) gen.Dao {
					return tx.Not(
						gen.Exists(
							dao.TorrentContent.Where(
								dao.TorrentContent.InfoHash.EqCol(
									dao.Torrent.InfoHash,
								),
							),
						),
					)
				})
			}

			go func() {
				err := handleJob(
					ctx,
					*msg,
					dao,
					scopes,
					processJobProvider,
					batchJobProvider,
				)

				if err == nil {
					err = runner.ErrCompleted
				}

				cancel(err)
			}()

			return runner.NopShutdowner, nil
		}
	}
}

func handleJob(
	ctx context.Context,
	msg batch.MessageParams,
	dao *dao.Query,
	scopes []func(gen.Dao) gen.Dao,
	processJobProvider queue.JobProvider[processor.MessageParams],
	batchJobProvider queue.JobProvider[batch.MessageParams],
) error {
	priority := 10
	// prioritise jobs where API calls are disabled as they will run faster:
	if msg.ApisDisabled() {
		priority = 4
	}

	maxInfoHash := msg.InfoHashGreaterThan
	chunkSize := uint(0)
	done := false

	var queueJobs []*model.QueueJob

	for {
		torrents, err := dao.Torrent.WithContext(ctx).
			Scopes(scopes...).
			Where(
				dao.Torrent.InfoHash.Gt(maxInfoHash),
				dao.Torrent.UpdatedAt.Lt(msg.UpdatedBefore),
			).
			Select(dao.Torrent.InfoHash).
			Order(dao.Torrent.InfoHash).
			Limit(int(msg.BatchSize)).
			Find()
		if err != nil {
			return err
		}

		if len(torrents) == 0 {
			done = true
			break
		}

		var infoHashes []protocol.ID

		for _, t := range torrents {
			maxInfoHash = t.InfoHash
			infoHashes = append(infoHashes, t.InfoHash)
			chunkSize++
		}

		job, err := processJobProvider(processor.MessageParams{
			ClassifierParams: processor.ClassifierParams{
				ClassifyMode:       msg.ClassifyMode,
				ClassifierWorkflow: msg.ClassifierWorkflow,
				ClassifierFlags:    msg.ClassifierFlags,
			},
			InfoHashes: infoHashes,
		}, model.QueueJobPriority(priority))
		if err != nil {
			return err
		}

		queueJobs = append(queueJobs, &job)

		if len(torrents) < int(msg.BatchSize) {
			done = true
			break
		}

		if chunkSize >= msg.ChunkSize {
			break
		}
	}

	if !done {
		job, err := batchJobProvider(batch.MessageParams{
			InfoHashGreaterThan: maxInfoHash,
			UpdatedBefore:       msg.UpdatedBefore,
			ClassifyMode:        msg.ClassifyMode,
			ClassifierWorkflow:  msg.ClassifierWorkflow,
			ClassifierFlags:     msg.ClassifierFlags,
			ChunkSize:           msg.ChunkSize,
			BatchSize:           msg.BatchSize,
			ContentTypes:        msg.ContentTypes,
			Orphans:             msg.Orphans,
		})
		if err != nil {
			return err
		}

		queueJobs = append(queueJobs, &job)
	}

	if len(queueJobs) > 0 {
		if err := dao.QueueJob.
			WithContext(ctx).
			Create(queueJobs...); err != nil {
			return err
		}
	}

	return nil
}
