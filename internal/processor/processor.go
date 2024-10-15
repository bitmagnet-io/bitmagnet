package processor

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"golang.org/x/sync/semaphore"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

type Processor interface {
	Process(ctx context.Context, params MessageParams) error
}

type processor struct {
	defaultWorkflow  string
	search           search.Search
	runner           classifier.Runner
	dao              *dao.Query
	blockingManager  blocking.Manager
	processSemaphore *semaphore.Weighted
	persistSemaphore *semaphore.Weighted
}

type MissingHashesError struct {
	InfoHashes []protocol.ID
}

func (e MissingHashesError) Error() string {
	return fmt.Sprintf("missing %d info hashes", len(e.InfoHashes))
}

func (c processor) Process(ctx context.Context, params MessageParams) error {
	if err := c.processSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer c.processSemaphore.Release(1)
	workflowName := params.ClassifierWorkflow
	if workflowName == "" {
		workflowName = c.defaultWorkflow
	}
	searchResult, searchErr := c.search.TorrentsWithMissingInfoHashes(
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
		return searchErr
	}
	tcResult, tcErr := c.search.TorrentContent(
		ctx,
		query.Where(search.TorrentContentInfoHashCriteria(params.InfoHashes...)),
		search.HydrateTorrentContentContent(),
	)
	if tcErr != nil {
		return tcErr
	}
	for _, tc := range tcResult.Items {
		for ti, t := range searchResult.Torrents {
			if t.InfoHash == tc.InfoHash {
				searchResult.Torrents[ti].Contents = append(searchResult.Torrents[ti].Contents, tc.TorrentContent)
				break
			}
		}
	}
	var errs []error
	failedHashes := make([]protocol.ID, 0, len(searchResult.MissingInfoHashes))
	failedHashes = append(failedHashes, searchResult.MissingInfoHashes...)
	if len(failedHashes) > 0 {
		errs = append(errs, MissingHashesError{InfoHashes: failedHashes})
	}
	tcs := make([]model.TorrentContent, 0, len(searchResult.Torrents))
	var idsToDelete []string
	var infoHashesToDelete []protocol.ID
	tagsToAdd := make(map[protocol.ID]map[string]struct{})
	for _, torrent := range searchResult.Torrents {
		thisDeleteIds := make(map[string]struct{}, len(torrent.Contents))
		foundMatch := false
		for _, tc := range torrent.Contents {
			thisDeleteIds[tc.ID] = struct{}{}
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
		cl, classifyErr := c.runner.Run(ctx, workflowName, params.ClassifierFlags, torrent)
		if classifyErr != nil {
			if errors.Is(classifyErr, classification.ErrDeleteTorrent) {
				infoHashesToDelete = append(infoHashesToDelete, torrent.InfoHash)
			} else {
				failedHashes = append(failedHashes, torrent.InfoHash)
				errs = append(errs, classifyErr)
			}
			continue
		}
		torrentContent := newTorrentContent(torrent, cl)
		tcId := torrentContent.InferID()
		for id := range thisDeleteIds {
			if id != tcId {
				idsToDelete = append(idsToDelete, id)
			}
		}
		tcs = append(tcs, torrentContent)
		if len(cl.Tags) > 0 {
			tagsToAdd[torrent.InfoHash] = cl.Tags
		}
	}
	if len(failedHashes) > 0 {
		if len(tcs) == 0 {
			return errors.Join(errs...)
		}
		republishJob, republishJobErr := NewQueueJob(MessageParams{
			InfoHashes:         failedHashes,
			ClassifyMode:       params.ClassifyMode,
			ClassifierWorkflow: workflowName,
			ClassifierFlags:    params.ClassifierFlags,
		})
		if republishJobErr != nil {
			return errors.Join(append(errs, republishJobErr)...)
		}
		if err := c.dao.QueueJob.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&republishJob); err != nil {
			return errors.Join(append(errs, err)...)
		}
	}
	if len(tcs) == 0 {
		return nil
	}
	if persistErr := c.persist(ctx, persistPayload{
		torrentContents:  tcs,
		deleteIds:        idsToDelete,
		deleteInfoHashes: infoHashesToDelete,
		addTags:          tagsToAdd,
	}); persistErr != nil {
		return persistErr
	}
	return nil
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
