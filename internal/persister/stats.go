package persister

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

func (s AllTablesStats) Add(tableName string, stats TableStats) {
	if s == nil {
		s = make(AllTablesStats)
	}
	existing := s[tableName]
	existing.Merge(stats)
	s[tableName] = existing
}
