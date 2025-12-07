package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/batch"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
)

func newPersistTorrentsWorker(
	classifierRunner classifier.Runner,
	persisterAdder persister.Adder,
	processorAdder batch.Adder[protocol.ID],
	scrapeAdder channel.Adder[nodeHasPeersForHash],
	size int,
	savePieces SavePieces,
	saveFilesThreshold int,
) channel.Worker[infoHashWithMetaInfo] {
	return channel.NewWorker(
		func(ctx context.Context, i infoHashWithMetaInfo) error {
			torrent := createTorrentModel(i.infoHash, i.metaInfo, savePieces, saveFilesThreshold)

			inputs := persister.Inputs{persister.InputTorrents(torrent)}

			if i.isVerifiedAbsentFromDB {
				classifyResult, err := classifierRunner.Run(ctx, "default", classifier.Flags{
					"apis_disabled":         true,
					"local_search_disabled": true,
				}, torrent)

				switch {
				case err == nil:
					inputs = append(inputs, persister.InputTorrentContent(classifyResult.ToTorrentContent()))
				case errors.Is(err, classification.ErrDeleteTorrent):
					inputs = persister.Inputs{persister.InputDeleteInfoHashes(i.infoHash)}
				default:
					return fmt.Errorf("failed to run classifier: %w", err)
				}
			}

			err := persisterAdder.Add(ctx, inputs.Input())
			if err != nil {
				return fmt.Errorf("failed to add torrent to persister: %w", err)
			}

			err = processorAdder.Add(ctx, torrent.InfoHash)
			if err != nil {
				return fmt.Errorf("failed to enqueue torrent processing: %w", err)
			}

			err = scrapeAdder.Add(ctx, i.nodeHasPeersForHash)
			if err != nil {
				return fmt.Errorf("failed to enqueue torrent scrape: %w", err)
			}

			return nil
		},
		channel.WithSize[infoHashWithMetaInfo](size),
	)
}

// runPersistTorrents waits on the persistTorrents channel, and persists torrents to the database in batches.
// After persisting each batch it will publish a message to the classifier,
// and forward the hash on the scrape channel to attempt finding the seeders/leechers.
// func (cr *crawler) runPersistTorrents(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
// 	shutdown := make(chan struct{})

// 	go func() {
// 		defer cancel(nil)

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return

// 			case <-shutdown:
// 				return

// 			case is := <-cr.classifyTorrents.Out():
// 				hashMap := make(map[protocol.ID]infoHashWithMetaInfo, len(is))

// 				var hashesToClassify []protocol.ID

// 				flushHashesToClassify := func() {
// 					if len(hashesToClassify) > 0 {
// 						job, err := cr.queueJobProvider(processor.MessageParams{
// 							InfoHashes: hashesToClassify,
// 						},
// 							// delay the classifier by a minute to allow time for the S/L
// 							// scrape:
// 							model.QueueJobDelayBy(time.Minute),
// 						)
// 						if err != nil {
// 							cr.logger.Errorf("error creating queue job: %s", err.Error())
// 						} else {
// 							err = cr.persister.Add(ctx, persister.InputQueueJobs(job))
// 							if err != nil {
// 								cr.logger.Errorf("error adding job to persister: %s", err.Error())
// 							}
// 						}
// 					}

// 					hashesToClassify = make([]protocol.ID, 0, classifyBatchSize)
// 				}
// 				flushHashesToClassify()

// 				var (
// 					mtx sync.Mutex
// 					wg  sync.WaitGroup
// 				)

// 				for _, i := range is {
// 					if _, ok := hashMap[i.infoHash]; ok {
// 						continue
// 					}

// 					hashMap[i.infoHash] = i

// 					wg.Add(1)

// 					go func() {
// 						defer wg.Done()

// 						torrent := createTorrentModel(i.infoHash, i.metaInfo, cr.savePieces, cr.saveFilesThreshold)

// 						inputs := persister.Inputs{persister.InputTorrents(torrent)}

// 						if i.isVerifiedAbsentFromDB {
// 							clResult, err := cr.classifier.Run(ctx, "default", classifier.Flags{
// 								"apis_disabled":         true,
// 								"local_search_disabled": true,
// 							}, torrent)
// 							switch {
// 							case err == nil:
// 								inputs = append(inputs, persister.InputTorrentContents(clResult.ToTorrentContent()))
// 							case errors.Is(err, classification.ErrDeleteTorrent):
// 								err = cr.blockingManager.Block(ctx, []protocol.ID{i.infoHash}, false)
// 								if err != nil {
// 									cr.logger.Errorf("error blocking torrent: %s", err.Error())
// 								}
// 								return
// 							default:
// 								cr.logger.Errorf("error running classifier: %s", err.Error())
// 							}
// 						}

// 						err := cr.persister.Add(ctx, inputs.Input())
// 						if err != nil {
// 							cr.logger.Errorf("error adding torrent to persister: %s", err.Error())
// 						}

// 						mtx.Lock()
// 						defer mtx.Unlock()

// 						hashesToClassify = append(hashesToClassify, i.infoHash)
// 						if len(hashesToClassify) >= classifyBatchSize {
// 							flushHashesToClassify()
// 						}
// 					}()
// 				}

// 				wg.Wait()

// 				flushHashesToClassify()
// 			}
// 		}
// 	}()

// 	return func(context.Context) error {
// 		close(shutdown)

// 		<-ctx.Done()

// 		return nil
// 	}, nil
// }

func createTorrentModel(
	hash protocol.ID,
	info metainfo.Info,
	savePieces SavePieces,
	saveFilesThreshold int,
) model.Torrent {
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

	files := make([]model.TorrentFile, 0, min(saveFilesThreshold, len(info.Files)))

	for i, file := range info.Files {
		if i >= saveFilesThreshold {
			filesStatus = model.FilesStatusOverThreshold
			break
		}

		// todo: Check this!
		var nullExt model.NullString

		if ext := strings.TrimPrefix(filepath.Ext(file.DisplayPath(&info)), "."); ext != "" {
			nullExt = model.NewNullString(ext)
		}

		files = append(files, model.TorrentFile{
			InfoHash:  hash,
			Index:     uint(i),
			Path:      file.DisplayPath(&info),
			Extension: nullExt,
			Size:      uint(file.Length),
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
		// CreatedAt: time.Now(),
	}
}
