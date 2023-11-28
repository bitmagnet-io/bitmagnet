package dhtcrawler

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"net/netip"
	"sync"
)

func (c *crawler) runRequestMetaInfo(ctx context.Context) {
	_ = c.requestMetaInfo.Run(ctx, func(req infoHashWithPeers) {
		mi, reqErr := c.doRequestMetaInfo(ctx, req.infoHash, req.peers)
		if reqErr != nil {
			return
		}
		select {
		case <-ctx.Done():
		case c.persistTorrents.In() <- infoHashWithMetaInfo{
			nodeHasPeersForHash: req.nodeHasPeersForHash,
			metaInfo:            mi.Info,
		}:
		}
	})
}

func (c *crawler) doRequestMetaInfo(
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
		res, err := c.metainfoRequester.Request(ctx, hash, p)
		if err != nil {
			addErr(err)
			continue
		}
		if banErr := c.banningChecker.Check(res.Info); banErr != nil {
			_ = c.blockingManager.Block(ctx, []protocol.ID{hash})
			return metainforequester.Response{}, banErr
		}
		return res, nil
	}
	return metainforequester.Response{}, errors.Join(errs...)
}
