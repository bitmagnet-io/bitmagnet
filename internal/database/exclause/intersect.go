package exclause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Intersect is intersect clause
type Intersect struct {
	Statements []clause.Expression
}

// Name intersect clause name
func (intersect Intersect) Name() string {
	return "INTERSECT"
}

// Build build intersect clause
func (intersect Intersect) Build(builder clause.Builder) {
	for index, statement := range intersect.Statements {
		if index != 0 {
			builder.WriteString(" INTERSECT ")
		}
		statement.Build(builder)
	}
}

// MergeClause merge Intersect clauses
func (intersect Intersect) MergeClause(mergeClause *clause.Clause) {
	if u, ok := mergeClause.Expression.(Intersect); ok {
		statements := make([]clause.Expression, len(u.Statements)+len(intersect.Statements))
		copy(statements, u.Statements)
		copy(statements[len(u.Statements):], intersect.Statements)
		intersect.Statements = statements
	}

	mergeClause.Expression = intersect
}

// NewIntersect is easy to create new Intersect
//
//	// examples
//	// SELECT * FROM `general_users` INTERSECT SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewIntersect("SELECT * FROM `admin_users`")).Scan(&users)
//
//	// SELECT * FROM `general_users` INTERSECT SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewIntersect(db.Table("admin_users"))).Scan(&users)
//
//	// SELECT * FROM `general_users` INTERSECT ALL SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewIntersect("ALL ?", db.Table("admin_users"))).Scan(&users)
func NewIntersect(query interface{}, args ...interface{}) Intersect {
	switch v := query.(type) {
	case *gorm.DB:
		return Intersect{
			Statements: []clause.Expression{Subquery{DB: v}},
		}
	case string:
		return Intersect{
			Statements: []clause.Expression{clause.Expr{SQL: v, Vars: args}},
		}
	}
	return Intersect{}
}
