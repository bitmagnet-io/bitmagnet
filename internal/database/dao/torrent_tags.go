package dao

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gorm/clause"
)

func (t *torrentTag) Put(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error {
	if len(infoHashes) == 0 || len(tagNames) == 0 {
		return nil
	}
	tagMap := make(map[string]struct{}, len(tagNames))
	for _, tagName := range tagNames {
		if validateErr := model.ValidateTagName(tagName); validateErr != nil {
			return validateErr
		}
		tagMap[tagName] = struct{}{}
	}
	hashMap := make(map[protocol.ID]struct{}, len(infoHashes))
	torrentTags := make([]*model.TorrentTag, 0, len(infoHashes)*len(tagMap))
	for _, infoHash := range infoHashes {
		if _, ok := hashMap[infoHash]; !ok {
			hashMap[infoHash] = struct{}{}
			for tagName := range tagMap {
				torrentTags = append(torrentTags, &model.TorrentTag{
					InfoHash: infoHash,
					Name:     tagName,
				})
			}
		}
	}
	return t.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(torrentTags, 100)
}

func (t *torrentTag) Set(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error {
	if len(infoHashes) == 0 {
		return nil
	}
	tagMap := make(map[string]struct{}, len(tagNames))
	for _, tagName := range tagNames {
		if validateErr := model.ValidateTagName(tagName); validateErr != nil {
			return validateErr
		}
		tagMap[tagName] = struct{}{}
	}
	hashMap := make(map[protocol.ID]struct{}, len(infoHashes))
	valuers := make([]driver.Valuer, 0, len(infoHashes))
	torrentTags := make([]*model.TorrentTag, 0, len(infoHashes)*len(tagMap))
	for _, infoHash := range infoHashes {
		if _, ok := hashMap[infoHash]; !ok {
			hashMap[infoHash] = struct{}{}
			valuers = append(valuers, infoHash)
			for tagName := range tagMap {
				torrentTags = append(torrentTags, &model.TorrentTag{
					InfoHash: infoHash,
					Name:     tagName,
				})
			}
		}
	}
	return Use(t.UnderlyingDB()).Transaction(func(tx *Query) error {
		if _, deleteErr := tx.TorrentTag.WithContext(ctx).Where(
			t.InfoHash.In(valuers...),
			t.Name.NotIn(tagNames...),
		).Delete(); deleteErr != nil {
			return deleteErr
		}
		return tx.TorrentTag.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).CreateInBatches(torrentTags, 100)
	})
}

func (t *torrentTag) Delete(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error {
	if len(infoHashes) == 0 && len(tagNames) == 0 {
		return nil
	}
	q := t.WithContext(ctx)
	if len(infoHashes) > 0 {
		valuers := make([]driver.Valuer, 0, len(infoHashes))
		for _, infoHash := range infoHashes {
			valuers = append(valuers, infoHash)
		}
		q = q.Where(t.InfoHash.In(valuers...))
	}
	if len(tagNames) > 0 {
		q = q.Where(t.Name.In(tagNames...))
	}
	_, err := q.Delete()
	return err
}
