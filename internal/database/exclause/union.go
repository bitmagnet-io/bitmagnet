package exclause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Union is union clause
type Union struct {
	Statements []clause.Expression
}

// Name union clause name
func (union Union) Name() string {
	return "UNION"
}

// Build build union clause
func (union Union) Build(builder clause.Builder) {
	for index, statement := range union.Statements {
		if index != 0 {
			builder.WriteString(" UNION ")
		}
		statement.Build(builder)
	}
}

// MergeClause merge Union clauses
func (union Union) MergeClause(mergeClause *clause.Clause) {
	if u, ok := mergeClause.Expression.(Union); ok {
		statements := make([]clause.Expression, len(u.Statements)+len(union.Statements))
		copy(statements, u.Statements)
		copy(statements[len(u.Statements):], union.Statements)
		union.Statements = statements
	}

	mergeClause.Expression = union
}

// NewUnion is easy to create new Union
//
//	// examples
//	// SELECT * FROM `general_users` UNION SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewUnion("SELECT * FROM `admin_users`")).Scan(&users)
//
//	// SELECT * FROM `general_users` UNION SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewUnion(db.Table("admin_users"))).Scan(&users)
//
//	// SELECT * FROM `general_users` UNION ALL SELECT * FROM `admin_users`
//	db.Table("general_users").Clauses(exclause.NewUnion("ALL ?", db.Table("admin_users"))).Scan(&users)
func NewUnion(query interface{}, args ...interface{}) Union {
	switch v := query.(type) {
	case *gorm.DB:
		return Union{
			Statements: []clause.Expression{Subquery{DB: v}},
		}
	case string:
		return Union{
			Statements: []clause.Expression{clause.Expr{SQL: v, Vars: args}},
		}
	}
	return Union{}
}
