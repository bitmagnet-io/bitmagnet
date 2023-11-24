package persistence

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gorm/clause"
	"strings"
)

type TorrentPersistence interface {
	PutTorrentTags(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error
	SetTorrentTags(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error
	DeleteTorrentTags(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error
}

func (p *persistence) PutTorrentTags(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error {
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
	return p.q.Transaction(func(tx *dao.Query) error {
		var existingHashes []protocol.ID
		if existingErr := tx.Torrent.WithContext(ctx).Where(
			tx.Torrent.InfoHash.In(valuers...),
		).Pluck(tx.Torrent.InfoHash, &existingHashes); existingErr != nil {
			return existingErr
		}
		var missingHashes []string
		for _, infoHash := range infoHashes {
			found := false
			for _, existingHash := range existingHashes {
				if infoHash == existingHash {
					found = true
					break
				}
			}
			if !found {
				missingHashes = append(missingHashes, infoHash.String())
			}
		}
		if len(missingHashes) > 0 {
			return errors.New("missing torrents: " + strings.Join(missingHashes, ", "))
		}
		return tx.TorrentTag.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).CreateInBatches(torrentTags, 100)
	})
}

func (p *persistence) SetTorrentTags(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error {
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
	return p.q.Transaction(func(tx *dao.Query) error {
		if _, deleteErr := tx.TorrentTag.WithContext(ctx).Where(
			tx.TorrentTag.InfoHash.In(valuers...),
			tx.TorrentTag.Name.NotIn(tagNames...),
		).Delete(); deleteErr != nil {
			return deleteErr
		}
		return tx.TorrentTag.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).CreateInBatches(torrentTags, 100)
	})
}

func (p *persistence) DeleteTorrentTags(ctx context.Context, infoHashes []protocol.ID, tagNames []string) error {
	if len(infoHashes) == 0 && len(tagNames) == 0 {
		return nil
	}
	q := p.q.TorrentTag.WithContext(ctx)
	if len(infoHashes) > 0 {
		valuers := make([]driver.Valuer, 0, len(infoHashes))
		for _, infoHash := range infoHashes {
			valuers = append(valuers, infoHash)
		}
		q = q.Where(p.q.TorrentTag.InfoHash.In(valuers...))
	}
	if len(tagNames) > 0 {
		q = q.Where(p.q.TorrentTag.Name.In(tagNames...))
	}
	_, err := q.Delete()
	return err
}
