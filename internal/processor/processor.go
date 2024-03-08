package processor

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/workflow"
	"golang.org/x/sync/semaphore"
	"gorm.io/gen/field"
)

type Processor interface {
	Process(ctx context.Context, params MessageParams) error
}

type processor struct {
	search   search.Search
	workflow workflow.Workflow
	//classifier       classifier.Classifier
	dao              *dao.Query
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
	searchResult, searchErr := c.search.TorrentsWithMissingInfoHashes(
		ctx,
		params.InfoHashes,
		query.Preload(func(q *dao.Query) []field.RelationField {
			return []field.RelationField{
				q.Torrent.Files.RelationField,
				q.Torrent.Hint.RelationField,
				q.Torrent.Contents.RelationField,
			}
		}),
	)
	if searchErr != nil {
		return searchErr
	}
	var errs []error
	var failedHashes []protocol.ID
	if len(failedHashes) > 0 {
		errs = append(errs, MissingHashesError{InfoHashes: failedHashes})
		failedHashes = append(failedHashes, searchResult.MissingInfoHashes...)
	}
	tcs := make([]model.TorrentContent, 0, len(searchResult.Torrents))
	var deleteIds []string
	var deleteInfoHashes []protocol.ID
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
		//if params.ClassifyMode == ClassifyModeSkipUnmatched && torrent.Hint.IsNil() {
		//	useClassifier = classifier.FallbackClassifier{}
		//}
		classification, classifyErr := c.workflow.Run(ctx, torrent)
		if classifyErr != nil {
			if errors.Is(classifyErr, workflow.ErrDeleteTorrent) {
				deleteInfoHashes = append(deleteInfoHashes, torrent.InfoHash)
			} else {
				errs = append(errs, classifyErr)
				failedHashes = append(failedHashes, torrent.InfoHash)
			}
			continue
		}
		torrentContent := newTorrentContent(torrent, classification)
		tcId := torrentContent.InferID()
		for id := range thisDeleteIds {
			if id != tcId {
				deleteIds = append(deleteIds, id)
			}
		}
		tcs = append(tcs, torrentContent)
	}
	if len(failedHashes) > 0 {
		if len(tcs) == 0 {
			return errors.Join(errs...)
		}
		republishJob, republishJobErr := NewQueueJob(MessageParams{
			InfoHashes:   failedHashes,
			ClassifyMode: params.ClassifyMode,
		})
		if republishJobErr != nil {
			return errors.Join(append(errs, republishJobErr)...)
		}
		if err := c.dao.QueueJob.WithContext(ctx).Create(&republishJob); err != nil {
			return errors.Join(append(errs, err)...)
		}
	}
	if len(tcs) == 0 {
		return nil
	}
	if persistErr := c.persist(ctx, persistPayload{
		torrentContents:  tcs,
		deleteIds:        deleteIds,
		deleteInfoHashes: deleteInfoHashes,
	}); persistErr != nil {
		return persistErr
	}
	return nil
}

func newTorrentContent(t model.Torrent, c classifier.Classification) model.TorrentContent {
	tc := model.TorrentContent{
		Torrent:         t,
		InfoHash:        t.InfoHash,
		ContentType:     c.ContentType,
		Languages:       c.Languages,
		Episodes:        c.Episodes,
		VideoResolution: c.VideoResolution,
		VideoSource:     c.VideoSource,
		VideoCodec:      c.VideoCodec,
		Video3d:         c.Video3d,
		VideoModifier:   c.VideoModifier,
		ReleaseGroup:    c.ReleaseGroup,
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
