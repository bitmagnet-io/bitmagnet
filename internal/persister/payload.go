package persister

import (
	"context"
	"database/sql/driver"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type hashWithID struct {
	hash protocol.ID
	id   string
}

type payload struct {
	shouldFlush            bool
	torrentSources         maps.InsertMap[string, model.TorrentSource]
	torrentsTorrentSources maps.InsertMap[hashWithID, model.TorrentsTorrentSource]
	torrentPieces          maps.InsertMap[protocol.ID, model.TorrentPieces]
	torrentFiles           maps.InsertMap[hashWithID, model.TorrentFile]
	torrents               maps.InsertMap[protocol.ID, model.Torrent]
	content                maps.InsertMap[model.ContentRef, model.Content]
	torrentContents        maps.InsertMap[model.TorrentContentRef, model.TorrentContent]
	deleteTorrentContent   maps.InsertMap[model.TorrentContentRef, struct{}]
	deleteInfoHashes       maps.InsertMap[protocol.ID, struct{}]
	torrentTags            maps.InsertMap[hashWithID, model.TorrentTag]
	queueJobs              maps.InsertMap[string, model.QueueJob]
}

func (p *payload) requiredInfoHashes() map[protocol.ID]struct{} {
	result := make(map[protocol.ID]struct{})

	for _, hid := range p.torrentsTorrentSources.Keys() {
		result[hid.hash] = struct{}{}
	}

	for _, hash := range p.torrentPieces.Keys() {
		result[hash] = struct{}{}
	}

	for _, hashID := range p.torrentFiles.Keys() {
		result[hashID.hash] = struct{}{}
	}

	for _, tcRef := range p.torrentContents.Keys() {
		result[tcRef.InfoHash] = struct{}{}
	}

	return result
}

func newPayload(payloads ...Input) *payload {
	payload := &payload{
		torrentSources:         maps.NewInsertMap[string, model.TorrentSource](),
		torrentsTorrentSources: maps.NewInsertMap[hashWithID, model.TorrentsTorrentSource](),
		torrentPieces:          maps.NewInsertMap[protocol.ID, model.TorrentPieces](),
		torrentFiles:           maps.NewInsertMap[hashWithID, model.TorrentFile](),
		torrents:               maps.NewInsertMap[protocol.ID, model.Torrent](),
		content:                maps.NewInsertMap[model.ContentRef, model.Content](),
		torrentContents:        maps.NewInsertMap[model.TorrentContentRef, model.TorrentContent](),
		torrentTags:            maps.NewInsertMap[hashWithID, model.TorrentTag](),
		deleteTorrentContent:   maps.NewInsertMap[model.TorrentContentRef, struct{}](),
		deleteInfoHashes:       maps.NewInsertMap[protocol.ID, struct{}](),
		//addTags:                maps.NewInsertMap[protocol.ID, maps.InsertMap[string, struct{}]](),
		queueJobs: maps.NewInsertMap[string, model.QueueJob](),
	}

	Inputs(payloads).Input()(payload)

	return payload
}

func (p *payload) len() int {
	return p.torrentSources.Len() +
		p.torrentsTorrentSources.Len() +
		p.torrentPieces.Len() +
		p.torrentFiles.Len() +
		p.torrents.Len() +
		p.content.Len() +
		p.torrentContents.Len() +
		p.torrentTags.Len() +
		p.deleteTorrentContent.Len() +
		p.deleteInfoHashes.Len() +
		p.queueJobs.Len()
}

func (p payload) missingInfoHashes(ctx context.Context, tx *dao.Query) (map[protocol.ID]struct{}, error) {
	requiredInfoHashes := p.requiredInfoHashes()

	if len(requiredInfoHashes) == 0 {
		return requiredInfoHashes, nil
	}

	valuers := make([]driver.Valuer, 0, len(requiredInfoHashes))
	for h := range requiredInfoHashes {
		valuers = append(valuers, h)
	}

	var result []*model.Torrent
	err := tx.Torrent.WithContext(ctx).Select(
		tx.Torrent.InfoHash,
	).Where(
		tx.Torrent.InfoHash.In(valuers...),
	).UnderlyingDB().Find(&result).Error
	if err != nil {
		return nil, err
	}

	for _, result := range result {
		delete(requiredInfoHashes, result.InfoHash)
	}

	return requiredInfoHashes, nil
}
