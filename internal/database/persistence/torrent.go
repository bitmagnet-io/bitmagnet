package persistence

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gorm/clause"
)

type TorrentPersistence interface {
	GetTorrent(ctx context.Context, infoHash protocol.ID) (torrent model.Torrent, err error)
	GetTorrents(ctx context.Context, infoHashes ...protocol.ID) (torrents []model.Torrent, missingInfoHashes []protocol.ID, err error)
	PutTorrent(ctx context.Context, torrent model.Torrent) error
	TorrentExists(ctx context.Context, infoHash protocol.ID) (bool, error)
	// GetPersistedInfoHashes returns the subset of provided hashes that are persisted in the database.
	GetPersistedInfoHashes(ctx context.Context, infoHashesToCheck []protocol.ID) ([]protocol.ID, error)
}

func (p *persistence) GetTorrent(ctx context.Context, infoHash protocol.ID) (t model.Torrent, _ error) {
	torrents, _, err := p.GetTorrents(ctx, infoHash)
	if err != nil {
		return t, err
	}
	if len(torrents) == 0 {
		return t, ErrRecordNotFound
	}
	return torrents[0], nil
}

func (p *persistence) GetTorrents(ctx context.Context, infoHashes ...protocol.ID) ([]model.Torrent, []protocol.ID, error) {
	valuers := make([]driver.Valuer, 0, len(infoHashes))
	for _, infoHash := range infoHashes {
		valuers = append(valuers, infoHash)
	}
	rawTorrents, findErr := p.q.WithContext(ctx).Torrent.Where(p.q.Torrent.InfoHash.In(valuers...)).Preload(
		p.q.Torrent.Files.RelationField.Order(
			p.q.TorrentFile.Index,
		),
		p.q.Torrent.Sources.RelationField,
		p.q.Torrent.Contents.RelationField.Order(
			p.q.TorrentContent.ContentType.IsNotNull(),
			p.q.TorrentContent.ContentID.IsNotNull(),
		),
	).Find()
	if findErr != nil {
		return nil, nil, findErr
	}
	torrents := make([]model.Torrent, 0, len(rawTorrents))
	missingInfoHashes := make([]protocol.ID, 0, len(infoHashes)-len(rawTorrents))
	foundInfoHashes := make(map[protocol.ID]struct{}, len(rawTorrents))
nextInfoHash:
	for _, h := range infoHashes {
		for _, t := range rawTorrents {
			if t.InfoHash == h {
				if _, ok := foundInfoHashes[h]; ok {
					continue nextInfoHash
				}
				foundInfoHashes[h] = struct{}{}
				torrents = append(torrents, *t)
				continue nextInfoHash
			}
		}
		missingInfoHashes = append(missingInfoHashes, h)
	}
	return torrents, missingInfoHashes, nil
}

func (p *persistence) PutTorrent(ctx context.Context, torrent model.Torrent) error {
	return p.q.WithContext(ctx).Torrent.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&torrent)
}

func (p *persistence) TorrentExists(ctx context.Context, infoHash protocol.ID) (bool, error) {
	count, err := p.q.WithContext(ctx).Torrent.Where(p.q.Torrent.InfoHash.Eq(infoHash)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (p *persistence) GetPersistedInfoHashes(ctx context.Context, infoHashesToCheck []protocol.ID) ([]protocol.ID, error) {
	valuers := make([]driver.Valuer, 0, len(infoHashesToCheck))
	for _, infoHash := range infoHashesToCheck {
		valuers = append(valuers, infoHash)
	}
	var persistedInfoHashes []protocol.ID
	if err := p.q.WithContext(ctx).Torrent.Where(p.q.Torrent.InfoHash.In(valuers...)).Pluck(p.q.Torrent.InfoHash, &persistedInfoHashes); err != nil {
		return nil, err
	}
	return persistedInfoHashes, nil
}
