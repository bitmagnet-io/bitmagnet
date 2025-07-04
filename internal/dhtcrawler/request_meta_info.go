package dhtcrawler

import (
	"context"
	"errors"
	"net/netip"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
)

func (cr *crawler) runRequestMetaInfo(ctx context.Context) error {
	return cr.requestMetaInfo.Run(ctx, func(req infoHashWithPeers) {
		mi, reqErr := cr.doRequestMetaInfo(ctx, req.infoHash, req.peers)
		if reqErr != nil {
			return
		}
		select {
		case <-ctx.Done():
		case cr.persistTorrents.In() <- infoHashWithMetaInfo{
			nodeHasPeersForHash: req.nodeHasPeersForHash,
			metaInfo:            mi.Info,
		}:
		}
	})
}

func (cr *crawler) doRequestMetaInfo(
	ctx context.Context,
	hash protocol.ID,
	peers []netip.AddrPort,
) (metainforequester.Response, error) {
	var errs []error

	errsMutex := sync.Mutex{}
	addErr := func(err error) {
		errsMutex.Lock()
		errs = append(errs, err)
		errsMutex.Unlock()
	}

	for _, p := range peers {
		res, err := cr.metainfoRequester.Request(ctx, hash, p)
		if err != nil {
			addErr(err)
			continue
		}

		if banErr := cr.banningChecker.Check(res.Info); banErr != nil {
			_ = cr.blockingManager.Block(ctx, []protocol.ID{hash}, false)
			return metainforequester.Response{}, banErr
		}

		return res, nil
	}

	return metainforequester.Response{}, errors.Join(errs...)
}
