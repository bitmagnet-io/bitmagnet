package exclause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Except is except clause
type Except struct {
	Statements []clause.Expression
}

// Name except clause name
func (except Except) Name() string {
	return "EXCEPT"
}

// Build build except clause
func (except Except) Build(builder clause.Builder) {
	for index, statement := range except.Statements {
		if index != 0 {
			builder.WriteString(" EXCEPT ")
		}
		statement.Build(builder)
	}
}

// MergeClause merge Except clauses
func (except Except) MergeClause(mergeClause *clause.Clause) {
	if u, ok := mergeClause.Expression.(Except); ok {
		statements := make([]clause.Expression, len(u.Statements)+len(except.Statements))
		copy(statements, u.Statements)
		copy(statements[len(u.Statements):], except.Statements)
		except.Statements = statements
	}

	mergeClause.Expression = except
}

// NewExcept is easy to create new Except
//
//	// examples
//	// SELECT * FROM `general_users` EXCEPT SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewExcept("SELECT * FROM `admin_users`")).Scan(&users)
//
//	// SELECT * FROM `general_users` EXCEPT SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewExcept(db.Table("admin_users"))).Scan(&users)
//
//	// SELECT * FROM `general_users` EXCEPT ALL SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewExcept("ALL ?", db.Table("admin_users"))).Scan(&users)
func NewExcept(query interface{}, args ...interface{}) Except {
	switch v := query.(type) {
	case *gorm.DB:
		return Except{
			Statements: []clause.Expression{Subquery{DB: v}},
		}
	case string:
		return Except{
			Statements: []clause.Expression{clause.Expr{SQL: v, Vars: args}},
		}
	}
	return Except{}
}
