package persister

import (
	"maps"
	"slices"

	"go.uber.org/zap"
)

type TableStats struct {
	Created  int
	Updated  int
	Deleted  int
	Affected int
	Ignored  int
}

func (s *TableStats) Merge(other TableStats) {
	s.Created += other.Created
	s.Updated += other.Updated
	s.Deleted += other.Deleted
	s.Affected += other.Affected
	s.Ignored += other.Ignored
}

type AllTablesStats map[string]TableStats

func (s AllTablesStats) LogSummary() []any {
	tableNames := slices.Collect(maps.Keys(s))
	slices.Sort(tableNames)

	var result []any
	for _, table := range tableNames {
		stats := s[table]
		if stats.Affected > 0 {
			result = append(result, table+":affected", stats.Affected)
		}
		if stats.Created > 0 {
			result = append(result, table+":created", stats.Created)
		}
		if stats.Updated > 0 {
			result = append(result, table+":updated", stats.Updated)
		}
		if stats.Deleted > 0 {
			result = append(result, table+":deleted", stats.Deleted)
		}
		if stats.Ignored > 0 {
			result = append(result, table+":ignored", stats.Ignored)
		}
	}

	return result
}

func (s AllTablesStats) LogFields() []zap.Field {
	tableNames := slices.Collect(maps.Keys(s))
	slices.Sort(tableNames)

	var result []zap.Field
	for _, table := range tableNames {
		stats := s[table]
		if stats.Affected > 0 {
			result = append(result, zap.Int(table+":affected", stats.Affected))
		}
		if stats.Created > 0 {
			result = append(result, zap.Int(table+":created", stats.Created))
		}
		if stats.Updated > 0 {
			result = append(result, zap.Int(table+":updated", stats.Updated))
		}
		if stats.Deleted > 0 {
			result = append(result, zap.Int(table+":deleted", stats.Deleted))
		}
		if stats.Ignored > 0 {
			result = append(result, zap.Int(table+":ignored", stats.Ignored))
		}
	}

	return result
}

func (s AllTablesStats) Add(tableName string, stats TableStats) {
	if s == nil {
		s = make(AllTablesStats)
	}
	existing := s[tableName]
	existing.Merge(stats)
	s[tableName] = existing
}
