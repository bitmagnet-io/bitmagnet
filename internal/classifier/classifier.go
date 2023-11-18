package classifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type Classifier interface {
	Classify(ctx context.Context, infoHashes ...protocol.ID) error
}

type classifier struct {
	resolver    Resolver
	dao         *dao.Query
	persistence persistence.Persistence
}

type MissingHashesError struct {
	InfoHashes []protocol.ID
}

func (e MissingHashesError) Error() string {
	return fmt.Sprintf("missing %d info hashes", len(e.InfoHashes))
}

func (c classifier) Classify(ctx context.Context, infoHashes ...protocol.ID) error {
	torrents, missingHashes, findErr := c.persistence.GetTorrents(ctx, infoHashes...)
	if findErr != nil {
		return findErr
	}
	resolved := make([]model.TorrentContent, 0, len(torrents))
	for _, torrent := range torrents {
		var torrentContent model.TorrentContent
		if len(torrent.Contents) > 0 {
			torrentContent = torrent.Contents[0]
			torrentContent.ContentSource = model.NullString{}
			torrentContent.ContentID = model.NullString{}
			torrentContent.Torrent = torrent
			torrentContent.Torrent.Contents = nil
		} else {
			torrentContent = model.TorrentContent{
				InfoHash: infoHashes[0],
				Torrent:  torrent,
			}
		}
		r, resolveErr := c.resolver.Resolve(ctx, torrentContent)
		if resolveErr != nil {
			if errors.Is(resolveErr, ErrNoMatch) {
				continue
			}
			return resolveErr
		}
		r.Torrent = model.Torrent{}
		resolved = append(resolved, r)
	}
	if resolveErr := c.Persist(ctx, resolved...); resolveErr != nil {
		return resolveErr
	}
	if len(missingHashes) > 0 {
		return MissingHashesError{
			InfoHashes: missingHashes,
		}
	}
	return nil
}
