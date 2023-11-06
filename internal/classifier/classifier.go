package classifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Persistence persistence.Persistence
	Resolver    resolver.RootResolver
	Logger      *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Classifier Classifier
}

type Classifier interface {
	Classify(ctx context.Context, infoHashes ...protocol.ID) error
}

func New(p Params) (Result, error) {
	return Result{
		Classifier: classifier{
			p.Persistence,
			p.Resolver,
			p.Logger,
		},
	}, nil
}

type classifier struct {
	p persistence.Persistence
	r resolver.RootResolver
	l *zap.SugaredLogger
}

type MissingHashesError struct {
	InfoHashes []protocol.ID
}

func (e MissingHashesError) Error() string {
	return fmt.Sprintf("missing %d info hashes", len(e.InfoHashes))
}

func (c classifier) Classify(ctx context.Context, infoHashes ...protocol.ID) error {
	torrents, missingHashes, findErr := c.p.GetTorrents(ctx, infoHashes...)
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
		r, resolveErr := c.r.Resolve(ctx, torrentContent)
		if resolveErr != nil {
			if errors.Is(resolveErr, resolver.ErrNoMatch) {
				continue
			}
			return resolveErr
		}
		r.Torrent = model.Torrent{}
		resolved = append(resolved, r)
	}
	if resolveErr := c.r.Persist(ctx, resolved...); resolveErr != nil {
		return resolveErr
	}
	if len(missingHashes) > 0 {
		return MissingHashesError{
			InfoHashes: missingHashes,
		}
	}
	return nil
}
