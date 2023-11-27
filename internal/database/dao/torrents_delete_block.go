package dao

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gorm"
)

const blockedTorrentsBloomFilterKey = "blocked_torrents"

func (q *Query) DeleteAndBlockTorrents(ctx context.Context, infoHashes []protocol.ID) (model.BloomFilter, error) {
	var valuers []driver.Valuer
	for _, infoHash := range infoHashes {
		valuers = append(valuers, infoHash)
	}
	var bf model.BloomFilter
	if txErr := q.Transaction(func(tx *Query) error {
		if _, deleteErr := tx.Torrent.WithContext(ctx).Where(tx.Torrent.InfoHash.In(valuers...)).Delete(); deleteErr != nil {
			return deleteErr
		}
		pBf, txErr := blockTx(ctx, tx, infoHashes)
		if txErr != nil {
			return txErr
		}
		bf = *pBf
		return nil
	}); txErr != nil {
		return bf, txErr
	}
	return bf, nil
}

func (q *Query) BlockTorrents(ctx context.Context, infoHashes []protocol.ID) (model.BloomFilter, error) {
	var bf model.BloomFilter
	if txErr := q.Transaction(func(tx *Query) error {
		pBf, txErr := blockTx(ctx, tx, infoHashes)
		if txErr != nil {
			return txErr
		}
		bf = *pBf
		return nil
	}); txErr != nil {
		return bf, txErr
	}
	return bf, nil
}

func blockTx(ctx context.Context, tx *Query, infoHashes []protocol.ID) (*model.BloomFilter, error) {
	bf, bfErr := tx.BloomFilter.WithContext(ctx).Where(
		tx.BloomFilter.Key.Eq(blockedTorrentsBloomFilterKey),
	).First()
	if errors.Is(bfErr, gorm.ErrRecordNotFound) {
		bf = &model.BloomFilter{
			Key:    blockedTorrentsBloomFilterKey,
			Filter: *bloom.NewDefaultStableBloomFilter(),
		}
	} else if bfErr != nil {
		return nil, bfErr
	}
	if len(infoHashes) == 0 {
		return bf, nil
	}
	for _, infoHash := range infoHashes {
		bf.Filter.Add(infoHash[:])
	}
	if saveErr := tx.BloomFilter.WithContext(ctx).Save(bf); saveErr != nil {
		return nil, saveErr
	}
	return bf, nil
}
