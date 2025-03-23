package exclause

import "gorm.io/gorm"

// ExtraClausePlugin support plugin that not supported clause by gorm
type ExtraClausePlugin struct{}

// Name return plugin name
func (*ExtraClausePlugin) Name() string {
	return "ExtraClausePlugin"
}

// Initialize register BuildClauses
func (*ExtraClausePlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Clauses = []string{
		"WITH", "SELECT", "FROM", "WHERE", "GROUP BY", "UNION", "INTERSECT", "EXCEPT", "ORDER BY", "LIMIT", "FOR",
	}
	db.Callback().Row().Clauses = []string{
		"WITH", "SELECT", "FROM", "WHERE", "GROUP BY", "UNION", "INTERSECT", "EXCEPT", "ORDER BY", "LIMIT", "FOR",
	}

	return nil
}

// New create new ExtraClausePlugin
//
//	// example
//	db.Use(extraClausePlugin.New())
func New() *ExtraClausePlugin {
	return &ExtraClausePlugin{}
}
