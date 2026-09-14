package dhtcrawler

import (
	"context"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gen/field"
)

// rescrapeStaleTickInterval is how often runRescrapeStale looks for a new batch of candidates.
// Kept short relative to rescrapeThreshold/confirmDeleteCooldown (which are usually many hours
// or days), so that a batch limit doesn't create a large multi-tick backlog on a DB with many
// stale torrents.
const rescrapeStaleTickInterval = 30 * time.Second

// rescrapeStaleBatchSize caps how many torrents are re-scraped per query per tick, so a DB with
// millions of stale torrents doesn't get walked in one enormous burst.
const rescrapeStaleBatchSize = 100

// rescrapeStaleConcurrency caps how many scrape requests runRescrapeStale has in flight at once.
const rescrapeStaleConcurrency = 10

// runRescrapeStale proactively re-checks torrents' seeders/leechers via a live DHT scrape,
// independent of whether the crawler happens to rediscover them again on its own - unlike
// infohash_triage.go's rescrapeThreshold check, which only fires reactively when a hash is
// re-observed. A torrent is only ever deleted after two independent scrapes,
// confirmDeleteCooldown apart, both come back below minSeeders: a single low reading only
// records the result and schedules a confirmation re-check, since a BEP-33 scrape result is
// only an approximation from whichever single DHT node answers, and can undercount a healthy
// swarm.
func (c *crawler) runRescrapeStale(ctx context.Context) {
	ticker := time.NewTicker(rescrapeStaleTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.rescrapeBatch(ctx, false)
			c.rescrapeBatch(ctx, true)
		}
	}
}

// rescrapeBatch finds and re-scrapes one batch of candidates. When confirming is false, it
// looks for torrents whose dht source hasn't been checked in over rescrapeThreshold and whose
// last-known seeders weren't already suspected low - a low result here is only a first strike,
// never deleted on this pass. When confirming is true, it looks specifically for torrents
// already suspected low (from a prior first-strike reading), re-checking them after the
// shorter confirmDeleteCooldown - a second consecutive low result here deletes the torrent.
func (c *crawler) rescrapeBatch(ctx context.Context, confirming bool) {
	var results []staleSource

	q := c.dao.TorrentsTorrentSource.WithContext(ctx).Select(
		c.dao.TorrentsTorrentSource.InfoHash,
	).Where(
		c.dao.TorrentsTorrentSource.Source.Eq("dht"),
	)

	if confirming {
		q = q.Where(
			c.dao.TorrentsTorrentSource.Seeders.Lt(model.NewNullUint(c.minSeeders)),
			c.dao.TorrentsTorrentSource.UpdatedAt.Lt(time.Now().Add(-c.confirmDeleteCooldown)),
		)
	} else {
		q = q.Where(
			c.dao.TorrentsTorrentSource.UpdatedAt.Lt(time.Now().Add(-c.rescrapeThreshold)),
			field.Or(
				c.dao.TorrentsTorrentSource.Seeders.IsNull(),
				c.dao.TorrentsTorrentSource.Seeders.Gte(model.NewNullUint(c.minSeeders)),
			),
		)
	}

	if err := q.Limit(rescrapeStaleBatchSize).Scan(&results); err != nil {
		c.logger.Errorf("error querying stale torrent sources: %s", err.Error())

		return
	}

	sem := make(chan struct{}, rescrapeStaleConcurrency)

	for _, r := range results {
		select {
		case <-ctx.Done():
			return
		case sem <- struct{}{}:
		}

		go func(infoHash protocol.ID) {
			defer func() { <-sem }()
			c.rescrapeAndCheck(ctx, infoHash, confirming)
		}(r.InfoHash)
	}

	// drain the semaphore so this batch's requests finish before the next tick starts a new one
	for range cap(sem) {
		select {
		case <-ctx.Done():
			return
		case sem <- struct{}{}:
		}
	}
}

type staleSource struct {
	InfoHash protocol.ID
}

// rescrapeAndCheck picks a node to ask, requests a fresh BEP-33 scrape for infoHash, and either
// deletes the torrent (if confirming a second consecutive low reading) or persists the result
// via the normal persistSources pipeline (matching every other scrape result, whether this is a
// routine re-check or a first-strike low reading awaiting confirmation).
func (c *crawler) rescrapeAndCheck(ctx context.Context, infoHash protocol.ID, confirming bool) {
	node, ok := pickNodeForHash(c.kTable.GetHashOrClosestNodes(infoHash))
	if !ok {
		return
	}

	result, err := c.requestScrape(ctx, nodeHasPeersForHash{infoHash: infoHash, node: node})
	if err != nil {
		return
	}

	if confirming && uint(result.bfsd.ApproximatedSize()) < c.minSeeders {
		if blockErr := c.blockingManager.Block(ctx, []protocol.ID{infoHash}, false); blockErr != nil {
			c.logger.Errorf("error blocking hash before delete: %s", blockErr.Error())

			return
		}

		if _, delErr := c.dao.Torrent.WithContext(ctx).Where(
			c.dao.Torrent.InfoHash.Eq(infoHash),
		).Delete(); delErr != nil {
			c.logger.Errorf("error deleting torrent with too few confirmed seeders: %s", delErr.Error())

			return
		}

		c.deletedTotal.With(prometheus.Labels{"reason": "min_seeders"}).Inc()

		return
	}

	select {
	case <-ctx.Done():
	case c.persistSources.In() <- result:
	}
}

// pickNodeForHash picks a node to query for a fresh scrape of a hash that may no longer be
// actively tracked in the k-table: prefer a node already known to have peers for the hash, and
// fall back to the closest known nodes by XOR distance (the same node selection a fresh DHT
// lookup would converge on).
func pickNodeForHash(res ktable.GetHashOrClosestNodesResult) (netip.AddrPort, bool) {
	if res.Found {
		if peers := res.Hash.Peers(); len(peers) > 0 {
			return peers[0].Addr, true
		}
	}

	if len(res.ClosestNodes) > 0 {
		return res.ClosestNodes[0].Addr(), true
	}

	return netip.AddrPort{}, false
}
