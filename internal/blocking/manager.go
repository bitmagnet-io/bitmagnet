package blocking

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"sync"
	"time"
)

type Manager interface {
	Filter(ctx context.Context, hashes []protocol.ID) ([]protocol.ID, error)
	Block(ctx context.Context, hashes []protocol.ID) error
	Flush(ctx context.Context) error
}

type manager struct {
	mutex         sync.Mutex
	dao           *dao.Query
	buffer        map[protocol.ID]struct{}
	filter        bloom.StableBloomFilter
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

func (m *manager) Block(ctx context.Context, hashes []protocol.ID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, hash := range hashes {
		m.buffer[hash] = struct{}{}
	}
	if m.shouldFlush() {
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

func (m *manager) flush(ctx context.Context) error {
	var hashes []protocol.ID
	for hash := range m.buffer {
		hashes = append(hashes, hash)
	}
	bf, blockErr := m.dao.BlockTorrents(ctx, hashes)
	if blockErr != nil {
		return blockErr
	}
	m.buffer = make(map[protocol.ID]struct{})
	m.filter = bf.Filter
	m.lastFlushedAt = time.Now()
	return nil
}

func (m *manager) shouldFlush() bool {
	return len(m.buffer) >= m.maxBufferSize || time.Since(m.lastFlushedAt) >= m.maxFlushWait
}
