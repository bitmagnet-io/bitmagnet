package search

import (
	"context"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSizeRangeCriteria_Raw(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		minBytes      *int64
		maxBytes      *int64
		key           string
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			name:          "no min or max specified",
			minBytes:      nil,
			maxBytes:      nil,
			key:           "torrent_contents.size",
			expectedQuery: "",
			expectedArgs:  []interface{}{},
		},
		{
			name:          "only min specified",
			minBytes:      int64Ptr(1024),
			maxBytes:      nil,
			key:           "torrent_contents.size",
			expectedQuery: "torrent_contents.size >= ?",
			expectedArgs:  []interface{}{int64(1024)},
		},
		{
			name:          "only max specified",
			minBytes:      nil,
			maxBytes:      int64Ptr(1048576),
			key:           "torrent_contents.size",
			expectedQuery: "torrent_contents.size <= ?",
			expectedArgs:  []interface{}{int64(1048576)},
		},
		{
			name:          "both min and max specified",
			minBytes:      int64Ptr(1024),
			maxBytes:      int64Ptr(1048576),
			key:           "torrent_contents.size",
			expectedQuery: "torrent_contents.size >= ? AND torrent_contents.size <= ?",
			expectedArgs:  []interface{}{int64(1024), int64(1048576)},
		},
		{
			name:          "different column name",
			minBytes:      int64Ptr(1024),
			maxBytes:      int64Ptr(1048576),
			key:           "torrents.size",
			expectedQuery: "torrents.size >= ? AND torrents.size <= ?",
			expectedArgs:  []interface{}{int64(1024), int64(1048576)},
		},
		{
			name:          "zero min value",
			minBytes:      int64Ptr(0),
			maxBytes:      int64Ptr(1048576),
			key:           "torrent_contents.size",
			expectedQuery: "torrent_contents.size >= ? AND torrent_contents.size <= ?",
			expectedArgs:  []interface{}{int64(0), int64(1048576)},
		},
		{
			name:          "min greater than max",
			minBytes:      int64Ptr(2048),
			maxBytes:      int64Ptr(1024),
			key:           "torrent_contents.size",
			expectedQuery: "torrent_contents.size >= ? AND torrent_contents.size <= ?",
			expectedArgs:  []interface{}{int64(2048), int64(1024)},
		},
		{
			name:          "large size values",
			minBytes:      int64Ptr(1_000_000_000), // 1GB
			maxBytes:      int64Ptr(1_000_000_000_000), // 1TB
			key:           "torrent_contents.size",
			expectedQuery: "torrent_contents.size >= ? AND torrent_contents.size <= ?",
			expectedArgs:  []interface{}{int64(1_000_000_000), int64(1_000_000_000_000)},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			criteria := SizeRangeCriteria{
				MinBytes: tt.minBytes,
				MaxBytes: tt.maxBytes,
				Key:      tt.key,
			}
			dbContext := &mockDBContext{}

			// Act
			rawCriteria, err := criteria.Raw(dbContext)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.expectedQuery, rawCriteria.Query)
			assert.Equal(t, tt.expectedArgs, rawCriteria.Args)
			assert.NotNil(t, rawCriteria.Joins)
			assert.Empty(t, rawCriteria.Joins.Entries(), "No joins should be required for size criteria")
		})
	}
}

// Helper function to create a pointer to an int64
func int64Ptr(v int64) *int64 {
	return &v
}

// Mock DbContext for testing Raw method
type mockDBContext struct{}

func (m *mockDBContext) Query() *dao.Query {
	return nil
}

func (m *mockDBContext) TableName() string {
	return ""
}

func (m *mockDBContext) NewSubQuery(ctx context.Context) query.SubQuery {
	return nil
}