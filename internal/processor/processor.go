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
	"golang.org/x/sync/semaphore"
	"gorm.io/gen/field"
)

type Processor interface {
	Process(ctx context.Context, params MessageParams) error
}

type processor struct {
	search           search.Search
	classifier       classifier.Classifier
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
	tcs := make([]model.TorrentContent, 0, len(searchResult.Torrents))
	for _, torrent := range searchResult.Torrents {
		if params.ClassifyMode != ClassifyModeRematch && !torrent.Hint.ContentSource.Valid {
			for _, tc := range torrent.Contents {
				if tc.ContentType.Valid &&
					tc.ContentSource.Valid &&
					(torrent.Hint.IsNil() || torrent.Hint.ContentType == tc.ContentType.ContentType) {
					torrent.Hint.ContentType = tc.ContentType.ContentType
					torrent.Hint.ContentSource = tc.ContentSource
					torrent.Hint.ContentID = tc.ContentID
					break
				}
			}
		}
		useClassifier := c.classifier
		if params.ClassifyMode == ClassifyModeSkipUnmatched && torrent.Hint.IsNil() {
			useClassifier = classifier.FallbackClassifier{}
		}
		classification, classifyErr := useClassifier.Classify(ctx, torrent)
		if classifyErr != nil {
			errs = append(errs, classifyErr)
			failedHashes = append(failedHashes, torrent.InfoHash)
			continue
		}
		torrentContent := newTorrentContent(torrent, classification)
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
	if persistErr := c.Persist(ctx, tcs...); persistErr != nil {
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
