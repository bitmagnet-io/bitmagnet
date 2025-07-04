package processor

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

const Namespace = "processor"

type Processor interface {
	NewJob(params MessageParams) runner.Runner
}

type processor struct {
	searchClient    search.Search
	daoProvider     database.DaoTransactionProvider
	blockingManager blocking.Blocker
	runner          classifier.Runner
	defaultWorkflow string
}

const concurrency = 10

func (c processor) NewJob(params MessageParams) runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		shutdown := make(chan struct{})

		setupCtx, setupCancel := context.WithCancel(ctx)

		go func() {
			// abort setup if shutdown initiated:
			select {
			case <-shutdown:
				cancel(fmt.Errorf("%w: %w: %w", Err, ErrSetup, ErrInterrupted))
			case <-setupCtx.Done():
			}

			setupCancel()
		}()

		searchResult, err := c.getTorrentsToClassify(setupCtx, params)

		setupCancel()

		if err != nil {
			err = fmt.Errorf("%w: %w: %w", Err, ErrSetup, err)
			cancel(err)

			return runner.NopShutdowner, err
		}

		if len(searchResult.Torrents)+len(searchResult.MissingInfoHashes) == 0 {
			cancel(runner.ErrCompleted)

			return runner.NopShutdowner, nil
		}

		go func() {
			defer cancel(runner.ErrCompleted)

			handleTorrentsResult, err := c.handleTorrents(
				ctx,
				shutdown,
				params.ClassifierParams,
				searchResult,
			)
			if err != nil {
				err = fmt.Errorf("%w: %w: %w", Err, ErrClassify, err)
				cancel(err)

				return
			}

			err = c.persist(ctx, handleTorrentsResult)
			if err != nil {
				err = fmt.Errorf("%w: %w: %w", Err, ErrPersist, err)

				cancel(err)
			}
		}()

		return func(shutdownCtx context.Context) error {
			close(shutdown)

			select {
			case <-ctx.Done():
			case <-shutdownCtx.Done():
			}

			cause := context.Cause(ctx)
			if cause != nil && !errors.Is(cause, runner.ErrCompleted) {
				return fmt.Errorf("%w: %w: %w", Err, ErrShutdown, cause)
			}

			return nil
		}, nil
	}
}

func (c processor) getTorrentsToClassify(
	ctx context.Context,
	params MessageParams,
) (search.TorrentsWithMissingInfoHashesResult, error) {
	searchResult, searchErr := c.searchClient.TorrentsWithMissingInfoHashes(
		ctx,
		params.InfoHashes,
		query.Preload(func(q *dao.Query) []field.RelationField {
			return []field.RelationField{
				q.Torrent.Files.RelationField,
				q.Torrent.Hint.RelationField,
				q.Torrent.Sources.RelationField,
			}
		}),
	)
	if searchErr != nil {
		return searchResult, searchErr
	}

	tcResult, tcErr := c.searchClient.TorrentContent(
		ctx,
		query.Where(search.TorrentContentInfoHashCriteria(params.InfoHashes...)),
		search.HydrateTorrentContentContent(),
	)
	if tcErr != nil {
		return searchResult, tcErr
	}

	for _, tc := range tcResult.Items {
		for ti, t := range searchResult.Torrents {
			if t.InfoHash == tc.InfoHash {
				searchResult.Torrents[ti].Contents = append(
					searchResult.Torrents[ti].Contents,
					tc.TorrentContent,
				)

				break
			}
		}
	}

	return searchResult, nil
}

func (c processor) handleTorrents(
	ctx context.Context,
	shutdown <-chan struct{},
	params ClassifierParams,
	searchResult search.TorrentsWithMissingInfoHashesResult,
) (persistPayload, error) {
	var (
		mtx                sync.Mutex
		idsToDelete        []string
		infoHashesToDelete []protocol.ID
		interrupted        bool
		errs               []error
	)

	tcs := make([]model.TorrentContent, 0, len(searchResult.Torrents))

	tagsToAdd := make(map[protocol.ID]map[string]struct{})

	failedHashes := make([]protocol.ID, 0, len(searchResult.MissingInfoHashes))
	failedHashes = append(failedHashes, searchResult.MissingInfoHashes...)

	sem := make(chan struct{}, concurrency)

	classifyCtx, classifyCancel := context.WithCancel(ctx)
	defer classifyCancel()

	for _, torrent := range searchResult.Torrents {
		select {
		case <-classifyCtx.Done():
			return persistPayload{}, classifyCtx.Err()
		case <-shutdown:
			mtx.Lock()
			failedHashes = append(failedHashes, torrent.InfoHash)
			interrupted = true
			mtx.Unlock()

			continue
		case sem <- struct{}{}:
		}

		go func(torrent model.Torrent) {
			defer func() {
				<-sem
			}()

			result, err := c.handleTorrent(classifyCtx, params, torrent)

			mtx.Lock()
			defer mtx.Unlock()

			if err != nil {
				if errors.Is(err, classification.ErrDeleteTorrent) {
					infoHashesToDelete = append(infoHashesToDelete, torrent.InfoHash)
				} else {
					failedHashes = append(failedHashes, torrent.InfoHash)
					errs = append(errs, err)
				}

				return
			}

			for id := range result.idsToDelete {
				idsToDelete = append(idsToDelete, id)
			}

			tcs = append(tcs, result.torrentContent)

			if len(result.tagsToAdd) > 0 {
				tagsToAdd[torrent.InfoHash] = result.tagsToAdd
			}
		}(torrent)
	}

	// wait for all classifications:
	for range concurrency {
		select {
		case <-ctx.Done():
			return persistPayload{}, ctx.Err()
		case sem <- struct{}{}:
		}
	}

	if len(tcs) == 0 && len(infoHashesToDelete) == 0 {
		err := ErrAllTorrentsFailed
		if interrupted {
			err = ErrInterrupted
		}

		err = fmt.Errorf("%w: %w", Err, errors.Join(append(errs, err)...))

		return persistPayload{}, err
	}

	return persistPayload{
		torrentContents:  tcs,
		deleteInfoHashes: infoHashesToDelete,
		deleteIDs:        idsToDelete,
		addTags:          tagsToAdd,
		failedInfoHashes: failedHashes,
		classifierParams: params,
	}, nil
}

type handleTorrentResult struct {
	torrentContent model.TorrentContent
	idsToDelete    map[string]struct{}
	tagsToAdd      map[string]struct{}
}

func (c processor) handleTorrent(
	ctx context.Context,
	params ClassifierParams,
	torrent model.Torrent,
) (handleTorrentResult, error) {
	thisDeleteIDs := make(map[string]struct{}, len(torrent.Contents))
	foundMatch := false

	for _, tc := range torrent.Contents {
		thisDeleteIDs[tc.ID] = struct{}{}

		if !foundMatch &&
			!torrent.Hint.ContentSource.Valid &&
			params.ClassifyMode != ClassifyModeRematch &&
			tc.ContentType.Valid &&
			tc.ContentSource.Valid &&
			(torrent.Hint.IsNil() || torrent.Hint.ContentType == tc.ContentType.ContentType) {
			torrent.Hint.ContentType = tc.ContentType.ContentType
			torrent.Hint.ContentSource = tc.ContentSource
			torrent.Hint.ContentID = tc.ContentID
			foundMatch = true
		}
	}

	workflowName := params.ClassifierWorkflow
	if workflowName == "" {
		workflowName = c.defaultWorkflow
	}

	cl, classifyErr := c.runner.Run(ctx, workflowName, params.ClassifierFlags, torrent)

	if classifyErr != nil {
		return handleTorrentResult{}, classifyErr
	}

	torrentContent := newTorrentContent(torrent, cl)

	delete(thisDeleteIDs, torrentContent.InferID())

	return handleTorrentResult{
		torrentContent: torrentContent,
		idsToDelete:    thisDeleteIDs,
		tagsToAdd:      cl.Tags,
	}, nil
}

func newTorrentContent(t model.Torrent, c classification.Result) model.TorrentContent {
	var filesCount model.NullUint
	if t.FilesCount.Valid {
		filesCount = t.FilesCount
	} else if t.FilesStatus == model.FilesStatusSingle {
		filesCount = model.NewNullUint(1)
	}

	tc := model.TorrentContent{
		Torrent:         t,
		InfoHash:        t.InfoHash,
		ContentType:     c.ContentType,
		Languages:       c.Languages,
		Episodes:        c.Episodes,
		VideoResolution: c.VideoResolution,
		VideoSource:     c.VideoSource,
		VideoCodec:      c.VideoCodec,
		Video3D:         c.Video3D,
		VideoModifier:   c.VideoModifier,
		ReleaseGroup:    c.ReleaseGroup,
		Size:            t.Size,
		FilesCount:      filesCount,
		Seeders:         t.Seeders(),
		Leechers:        t.Leechers(),
		PublishedAt:     t.PublishedAt(),
	}

	if c.Content != nil {
		content := *c.Content
		content.UpdateTsv()
		tc.ContentType = model.NewNullContentType(content.Type)
		tc.ContentSource = model.NewNullString(content.Source)
		tc.ContentID = model.NewNullString(content.ID)
		tc.Content = content
	}

	tc.UpdateTsv()

	return tc
}

type persistPayload struct {
	torrentContents  []model.TorrentContent
	deleteIDs        []string
	deleteInfoHashes []protocol.ID
	addTags          map[protocol.ID]map[string]struct{}
	failedInfoHashes []protocol.ID
	classifierParams ClassifierParams
}

func (c processor) persist(ctx context.Context, payload persistPayload) error {
	contentsMap := make(map[model.ContentRef]struct{}, len(payload.torrentContents))
	contentsPtr := make([]*model.Content, 0, len(payload.torrentContents))
	torrentContentsPtr := make([]*model.TorrentContent, 0, len(payload.torrentContents))
	torrentTagsPtr := make([]*model.TorrentTag, 0, len(payload.addTags))

	for _, tc := range payload.torrentContents {
		tcCopy := tc
		tcCopy.Torrent = model.Torrent{}

		if tcCopy.ContentID.Valid && tcCopy.Content.CreatedAt.IsZero() {
			contentRef := tcCopy.Content.Ref()
			if _, ok := contentsMap[contentRef]; !ok {
				contentsMap[contentRef] = struct{}{}
				contentCopy := tcCopy.Content
				contentsPtr = append(contentsPtr, &contentCopy)
			}
		}

		tcCopy.Content = model.Content{}
		torrentContentsPtr = append(torrentContentsPtr, &tcCopy)
	}

	for infoHash, tags := range payload.addTags {
		for tag := range tags {
			torrentTagsPtr = append(torrentTagsPtr, &model.TorrentTag{
				InfoHash: infoHash,
				Name:     tag,
			})
		}
	}

	if len(payload.deleteInfoHashes) > 0 {
		if blockErr := c.blockingManager.Block(ctx, payload.deleteInfoHashes, false); blockErr != nil {
			return blockErr
		}
	}

	return c.daoProvider.DaoTransaction(func(tx *dao.Query) error {
		if len(payload.failedInfoHashes) > 0 {
			republishJob, err := NewQueueJob(MessageParams{
				InfoHashes:       payload.failedInfoHashes,
				ClassifierParams: payload.classifierParams,
			})
			if err != nil {
				return err
			}

			if err := tx.QueueJob.WithContext(ctx).Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&republishJob); err != nil {
				return err
			}
		}

		if len(contentsPtr) > 0 {
			if createContentErr := tx.Content.WithContext(ctx).Clauses(
				clause.OnConflict{
					UpdateAll: true,
				}).CreateInBatches(contentsPtr, 100); createContentErr != nil {
				return createContentErr
			}
		}

		if len(payload.deleteIDs) > 0 {
			if _, deleteErr := tx.TorrentContent.WithContext(ctx).Where(
				tx.TorrentContent.ID.In(payload.deleteIDs...),
			).Delete(); deleteErr != nil {
				return deleteErr
			}
		}

		if len(torrentContentsPtr) > 0 {
			if createErr := tx.TorrentContent.WithContext(ctx).Clauses(
				clause.OnConflict{
					UpdateAll: true,
				},
			).CreateInBatches(torrentContentsPtr, 100); createErr != nil {
				return createErr
			}
		}

		if len(torrentTagsPtr) > 0 {
			if createErr := tx.TorrentTag.WithContext(ctx).Clauses(
				clause.OnConflict{
					DoNothing: true,
				},
			).CreateInBatches(torrentTagsPtr, 100); createErr != nil {
				return createErr
			}
		}

		if len(payload.deleteInfoHashes) > 0 {
			valuers := slice.Map(payload.deleteInfoHashes, func(infoHash protocol.ID) driver.Valuer {
				return infoHash
			})

			if _, deleteErr := tx.Torrent.WithContext(ctx).Where(
				tx.Torrent.InfoHash.In(valuers...),
			).Delete(); deleteErr != nil {
				return deleteErr
			}
		}

		return nil
	})
}
