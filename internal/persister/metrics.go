package persister

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func newFlusherMetrics(
	flusher iflusher,
	component *metrics.Component,
) *flusherMetrics {
	atm := make(allTablesMetrics, len(tableNames))

	for _, table := range tableNames {
		sub := component.MustSub(table)

		atm[table] = &tableMetrics{
			created:  sub.MustNewCounter("created"),
			updated:  sub.MustNewCounter("updated"),
			deleted:  sub.MustNewCounter("deleted"),
			affected: sub.MustNewCounter("affected"),
			ignored:  sub.MustNewCounter("ignored"),
		}
	}

	return &flusherMetrics{
		iflusher:         flusher,
		allTablesMetrics: atm,
		latencySampler:   component.MustNewSampler("latency"),
	}
}

type flusherMetrics struct {
	iflusher
	allTablesMetrics
	latencySampler *metrics.Sampler
}

func (m *flusherMetrics) flush(ctx context.Context, payload *payload) (AllTablesStats, error) {
	startTime := time.Now()

	allStats, err := m.iflusher.flush(ctx, payload)
	if err != nil {
		return allStats, err
	}

	m.latencySampler.Add(float32(time.Since(startTime).Seconds()))

	for table, stats := range allStats {
		metrics := m.allTablesMetrics[table]

		metrics.created.IncrN(stats.Created)
		metrics.updated.IncrN(stats.Updated)
		metrics.deleted.IncrN(stats.Deleted)
		metrics.affected.IncrN(stats.Affected)
		metrics.ignored.IncrN(stats.Ignored)
	}

	return allStats, err
}

type tableMetrics struct {
	created  *metrics.Counter
	updated  *metrics.Counter
	deleted  *metrics.Counter
	affected *metrics.Counter
	ignored  *metrics.Counter
}

type allTablesMetrics map[string]*tableMetrics

var tableNames = []string{
	model.TableNameTorrent,
	model.TableNameTorrentPieces,
	model.TableNameTorrentSource,
	model.TableNameTorrentsTorrentSource,
	model.TableNameContent,
	model.TableNameQueueJob,
	model.TableNameTorrentContent,
	model.TableNameTorrentFile,
	model.TableNameTorrentTag,
}
