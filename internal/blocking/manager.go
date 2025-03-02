package blocking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
	"time"
)

type Manager interface {
	Filter(ctx context.Context, hashes []protocol.ID) ([]protocol.ID, error)
	Block(ctx context.Context, hashes []protocol.ID, flush bool) error
	Flush(ctx context.Context) error
}

type manager struct {
	mutex         sync.Mutex
	pool          *pgxpool.Pool
	buffer        map[protocol.ID]struct{}
	filter        *bloom.StableBloomFilter
	maxBufferSize int
	lastFlushedAt time.Time
	maxFlushWait  time.Duration
}

func (m *manager) Filter(ctx context.Context, hashes []protocol.ID) ([]protocol.ID, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.filter.Cells() == 0 || m.shouldFlush() {
		if flushErr := m.flush(ctx); flushErr != nil {
			return nil, flushErr
		}
	}
	var filtered []protocol.ID
	for _, hash := range hashes {
		if _, ok := m.buffer[hash]; ok {
			continue
		}
		if m.filter.Test(hash[:]) {
			continue
		}
		filtered = append(filtered, hash)
	}
	return filtered, nil
}

func (m *manager) Block(ctx context.Context, hashes []protocol.ID, flush bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, hash := range hashes {
		m.buffer[hash] = struct{}{}
	}
	if flush || m.shouldFlush() {
		if flushErr := m.flush(ctx); flushErr != nil {
			return flushErr
		}
	}
	return nil
}

func (m *manager) Flush(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if len(m.buffer) == 0 {
		return nil
	}
	return m.flush(ctx)
}

const blockedTorrentsBloomFilterKey = "blocked_torrents"

func (m *manager) flush(ctx context.Context) error {
	var hashes []protocol.ID
	for hash := range m.buffer {
		hashes = append(hashes, hash)
	}

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, "delete from torrents where info_hash = any($1)", hashes)
	if err != nil {
		return fmt.Errorf("error deleting from torrents table: %w", err)
	}

	bf := bloom.NewDefaultStableBloomFilter()

	lobs := tx.LargeObjects()

	found := false

	var oid uint32
	var nullOid sql.NullInt32
	err = tx.QueryRow(ctx, "SELECT oid FROM bloom_filters WHERE key = $1", blockedTorrentsBloomFilterKey).Scan(&nullOid)
	if err == nil {
		found = true
		if nullOid.Valid {
			oid = uint32(nullOid.Int32)
			obj, err := lobs.Open(ctx, oid, pgx.LargeObjectModeRead)
			if err != nil {
				return fmt.Errorf("error opening large object for reading: %w", err)
			}
			_, err = bf.ReadFrom(obj)
			obj.Close()
			if err != nil {
				return fmt.Errorf("error reading current bloom filter: %w", err)
			}
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("error getting bloom filter object ID: %w", err)
	}

	if oid == 0 {
		// Create a new Large Object.
		// We pass 0, so the DB can pick an available oid for us.
		oid, err = lobs.Create(ctx, 0)
		if err != nil {
			return fmt.Errorf("error creating large object: %w", err)
		}
	}

	for _, hash := range hashes {
		bf.Add(hash[:])
	}

	obj, err := lobs.Open(ctx, oid, pgx.LargeObjectModeWrite)
	if err != nil {
		return fmt.Errorf("error opening large object for writing: %w", err)
	}

	_, err = bf.WriteTo(obj)
	if err != nil {
		return fmt.Errorf("error writing to large object: %w", err)
	}

	now := time.Now()
	if !found {
		_, err = tx.Exec(ctx, "INSERT INTO bloom_filters (key, oid, created_at, updated_at) VALUES ($1, $2, $3, $4)", blockedTorrentsBloomFilterKey, oid, now, now)
		if err != nil {
			return fmt.Errorf("error saving new bloom filter record: %w", err)
		}
	} else if !nullOid.Valid {
		_, err = tx.Exec(ctx, "UPDATE bloom_filters SET oid = $1, updated_at = $2 WHERE key = $3", oid, now, blockedTorrentsBloomFilterKey)
		if err != nil {
			return fmt.Errorf("error updating bloom filter record: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	m.buffer = make(map[protocol.ID]struct{})
	m.filter = bf
	m.lastFlushedAt = now

	return nil
}

func (m *manager) shouldFlush() bool {
	return len(m.buffer) >= m.maxBufferSize || time.Since(m.lastFlushedAt) >= m.maxFlushWait
}
