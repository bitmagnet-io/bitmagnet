package dhtcrawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm/clause"
)

// runPersistTorrents waits on the persistTorrents channel, and persists torrents to the database in batches.
// After persisting each batch it will publish a message to the classifier,
// and forward the hash on the scrape channel to attempt finding the seeders/leechers.
func (cr *crawler) runPersistTorrents(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
	shutdown := make(chan struct{})

	go func() {
		defer cancel(nil)

		for {
			select {
			case <-ctx.Done():
				return

			case <-shutdown:
				return

			case is := <-cr.persistTorrents.Out():
				torrentsToPersist := make([]*model.Torrent, 0, len(is))

				var torrentFilesToPersist []*model.TorrentFile

				var torrentSourcesToPersist []*model.TorrentsTorrentSource

				var torrentPiecesToPersist []*model.TorrentPieces

				var queueJobsToPersist []*model.QueueJob

				hashMap := make(map[protocol.ID]infoHashWithMetaInfo, len(is))

				var hashesToClassify []protocol.ID

				flushHashesToClassify := func() {
					if len(hashesToClassify) > 0 {
						job, err := processor.NewQueueJob(processor.MessageParams{
							InfoHashes: hashesToClassify,
						},
							// delay the classifier by a minute to allow time for the S/L
							// scrape:
							model.QueueJobDelayBy(time.Minute),
						)
						if err != nil {
							cr.logger.Errorf("error creating queue job: %s", err.Error())
						} else {
							queueJobsToPersist = append(queueJobsToPersist, &job)
						}
					}

					hashesToClassify = make([]protocol.ID, 0, classifyBatchSize)
				}
				flushHashesToClassify()

				for _, i := range is {
					if _, ok := hashMap[i.infoHash]; ok {
						continue
					}

					hashMap[i.infoHash] = i

					if t, err := createTorrentModel(
						i.infoHash, i.metaInfo, cr.savePieces, cr.saveFilesThreshold); err != nil {
						cr.logger.Errorf("error creating torrent model: %s", err.Error())
					} else {
						for _, f := range t.Files {
							fc := f
							torrentFilesToPersist = append(torrentFilesToPersist, &fc)
						}

						t.Files = nil
						for _, s := range t.Sources {
							sc := s
							torrentSourcesToPersist = append(torrentSourcesToPersist, &sc)
						}

						t.Sources = nil
						if cr.savePieces {
							pc := t.Pieces
							torrentPiecesToPersist = append(torrentPiecesToPersist, &pc)
							t.Pieces = model.TorrentPieces{}
						}

						torrentsToPersist = append(torrentsToPersist, &t)

						hashesToClassify = append(hashesToClassify, i.infoHash)
						if len(hashesToClassify) >= classifyBatchSize {
							flushHashesToClassify()
						}
					}
				}

				flushHashesToClassify()

				if persistErr := cr.daoProvider.DaoTransaction(func(tx *dao.Query) error {
					if err := tx.WithContext(ctx).Torrent.Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: string(tx.Torrent.InfoHash.ColumnName())}},
						DoUpdates: clause.AssignmentColumns([]string{
							string(tx.Torrent.Name.ColumnName()),
							string(tx.Torrent.FilesStatus.ColumnName()),
							string(tx.Torrent.FilesCount.ColumnName()),
							string(tx.Torrent.UpdatedAt.ColumnName()),
						}),
					}).CreateInBatches(torrentsToPersist, 100); err != nil {
						return err
					}
					if len(torrentFilesToPersist) > 0 {
						if err := tx.WithContext(ctx).TorrentFile.Clauses(clause.OnConflict{
							DoNothing: true,
						}).CreateInBatches(torrentFilesToPersist, 100); err != nil {
							return err
						}
					}
					if err := tx.WithContext(ctx).TorrentsTorrentSource.Clauses(clause.OnConflict{
						DoNothing: true,
					}).CreateInBatches(torrentSourcesToPersist, 100); err != nil {
						return err
					}
					if cr.savePieces {
						if err := tx.WithContext(ctx).TorrentPieces.Clauses(clause.OnConflict{
							DoNothing: true,
						}).CreateInBatches(torrentPiecesToPersist, 10); err != nil {
							return err
						}
					}

					return tx.WithContext(ctx).QueueJob.CreateInBatches(queueJobsToPersist, 10)
				}); persistErr != nil {
					cr.logger.Errorf("error persisting torrents: %s", persistErr)
				} else {
					cr.persistedTotal.With(prometheus.Labels{"entity": "Torrent"}).Add(float64(len(torrentsToPersist)))
					cr.logger.Debugw("persisted torrents", "count", len(torrentsToPersist))

					for _, i := range hashMap {
						select {
						case <-ctx.Done():
							return

						case <-shutdown:
							return

						case cr.scrape.In() <- i.nodeHasPeersForHash:
							continue
						}
					}
				}
			}
		}
	}()

	return func(context.Context) error {
		close(shutdown)

		<-ctx.Done()

		return nil
	}, nil
}

func createTorrentModel(
	hash protocol.ID,
	info metainfo.Info,
	savePieces bool,
	saveFilesThreshold uint,
) (model.Torrent, error) {
	name := info.BestName()

	private := false
	if info.Private != nil {
		private = *info.Private
	}

	var filesCount model.NullUint

	filesStatus := model.FilesStatusSingle
	if len(info.Files) > 0 {
		filesStatus = model.FilesStatusMulti
		filesCount = model.NewNullUint(uint(len(info.Files)))
	}

	files := make([]model.TorrentFile, 0, min(int(saveFilesThreshold), len(info.Files)))

	for i, file := range info.Files {
		if i >= int(saveFilesThreshold) {
			filesStatus = model.FilesStatusOverThreshold
			break
		}

		files = append(files, model.TorrentFile{
			InfoHash: hash,
			Index:    uint(i),
			Path:     file.DisplayPath(&info),
			Size:     uint(file.Length),
		})
	}

	var pieces model.TorrentPieces
	if savePieces {
		pieces = model.TorrentPieces{
			InfoHash:    hash,
			PieceLength: info.PieceLength,
			Pieces:      info.Pieces,
		}
	}

	return model.Torrent{
		InfoHash:    hash,
		Name:        name,
		Size:        uint(info.TotalLength()),
		Private:     private,
		Pieces:      pieces,
		Files:       files,
		FilesStatus: filesStatus,
		FilesCount:  filesCount,
		Sources: []model.TorrentsTorrentSource{
			{
				Source:   "dht",
				InfoHash: hash,
			},
		},
	}, nil
}

const classifyBatchSize = 100
