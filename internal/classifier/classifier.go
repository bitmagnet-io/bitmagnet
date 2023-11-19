package classifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type Classifier interface {
	Classify(ctx context.Context, infoHashes ...protocol.ID) error
}

type classifier struct {
	search   search.Search
	resolver Resolver
	dao      *dao.Query
}

type MissingHashesError struct {
	InfoHashes []protocol.ID
}

func (e MissingHashesError) Error() string {
	return fmt.Sprintf("missing %d info hashes", len(e.InfoHashes))
}

func (c classifier) Classify(ctx context.Context, infoHashes ...protocol.ID) error {
	searchResult, searchErr := c.search.TorrentsWithMissingInfoHashes(ctx, infoHashes)
	if searchErr != nil {
		return searchErr
	}
	resolved := make([]model.TorrentContent, 0, len(searchResult.Torrents))
	for _, torrent := range searchResult.Torrents {
		var torrentContent model.TorrentContent
		if len(torrent.Contents) > 0 {
			torrentContent = torrent.Contents[0]
			torrentContent.ContentSource = model.NullString{}
			torrentContent.ContentID = model.NullString{}
			torrentContent.Torrent = torrent
			torrentContent.Torrent.Contents = nil
		} else {
			torrentContent = model.TorrentContent{
				InfoHash: torrent.InfoHash,
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
	if len(searchResult.MissingInfoHashes) > 0 {
		return MissingHashesError{
			InfoHashes: searchResult.MissingInfoHashes,
		}
	}
	return nil
}
