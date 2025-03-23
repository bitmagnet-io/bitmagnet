package exclause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// With with clause
//
//		// examples
//		// WITH `cte` AS (SELECT * FROM `users`) SELECT * FROM `cte`
//		db.Clauses(exclause.With{CTEs: []exclause.CTE{{Name: "cte",
//	 	Subquery: clause.Expr{SQL: "SELECT * FROM `users`"}}}}).Table("cte").Scan(&users)
//
//		// WITH `cte` AS (SELECT * FROM `users`) SELECT * FROM `cte`
//		db.Clauses(exclause.With{CTEs: []exclause.CTE{{Name: "cte",
//	 	Subquery: exclause.Subquery{DB: db.Table("users")}}}}).Table("cte").Scan(&users)
//
//		// WITH `cte` (`id`,`name`) AS (SELECT * FROM `users`) SELECT * FROM `cte`
//		db.Clauses(exclause.With{CTEs: []exclause.CTE{{Name: "cte", Columns: []string{"id", "name"},
//	 	Subquery: exclause.Subquery{DB: db.Table("users")}}}}).Table("cte").Scan(&users)
//
//		// WITH RECURSIVE `cte` AS (SELECT * FROM `users`) SELECT * FROM `cte`
//		db.Clauses(exclause.With{Recursive: true, CTEs: []exclause.CTE{{Name: "cte",
//	 	Subquery: exclause.Subquery{DB: db.Table("users")}}}}).Table("cte").Scan(&users)
type With struct {
	Recursive    bool
	Materialized bool
	CTEs         []CTE
}

// CTE common table expressions
type CTE struct {
	Name     string
	Columns  []string
	Subquery clause.Expression
}

// Name with clause name
func (With) Name() string {
	return "WITH"
}

// Build build with clause
func (with With) Build(builder clause.Builder) {
	if with.Recursive {
		if _, err := builder.WriteString("RECURSIVE "); err != nil {
			panic(err)
		}
	}

	for index, cte := range with.CTEs {
		if index > 0 {
			if err := builder.WriteByte(','); err != nil {
				panic(err)
			}
		}

		cte.Build(builder, with.Materialized)
	}
}

// Build build CTE
func (cte CTE) Build(builder clause.Builder, materialized bool) {
	builder.WriteQuoted(cte.Name)

	if len(cte.Columns) > 0 {
		if _, err := builder.WriteString(" ("); err != nil {
			panic(err)
		}

		for index, column := range cte.Columns {
			if index > 0 {
				if err := builder.WriteByte(','); err != nil {
					panic(err)
				}
			}

			builder.WriteQuoted(column)
		}

		if err := builder.WriteByte(')'); err != nil {
			panic(err)
		}
	}

	if _, err := builder.WriteString(" AS "); err != nil {
		panic(err)
	}

	// Latest versions of Postgres default to non-materialized CTEs, so we don't need to
	// specify it explicitly. Sometimes you want to keep the optimisation fence though, in
	// which case you can set the Materialized flag to true.
	if materialized {
		if _, err := builder.WriteString("MATERIALIZED "); err != nil {
			panic(err)
		}
	}

	if err := builder.WriteByte('('); err != nil {
		panic(err)
	}

	cte.Subquery.Build(builder)

	if err := builder.WriteByte(')'); err != nil {
		panic(err)
	}
}

// MergeClause merge With clauses
func (with With) MergeClause(clause *clause.Clause) {
	if w, ok := clause.Expression.(With); ok {
		if w.Recursive {
			with.Recursive = true
		}

		ctes := make([]CTE, len(w.CTEs)+len(with.CTEs))
		copy(ctes, w.CTEs)
		copy(ctes[len(w.CTEs):], with.CTEs)
		with.CTEs = ctes
	}

	clause.Expression = with
}

// NewWith is easy to create new With
//
//		// examples
//		// WITH `cte` AS (SELECT * FROM `users`) SELECT * FROM `cte`
//		db.Clauses(exclause.NewWith("cte", "SELECT * FROM `users`")).Table("cte").Scan(&users)
//
//		// WITH `cte` AS (SELECT * FROM `users` WHERE `name` = 'WinterYukky') SELECT * FROM `cte`
//		db.Clauses(exclause.NewWith("cte", "SELECT * FROM `users` WHERE `name` = ?", "WinterYukky")).
//		  Table("cte").Scan(&users)
//
//		// WITH `cte` AS (SELECT * FROM `users` WHERE `name` = 'WinterYukky') SELECT * FROM `cte`
//		db.Clauses(exclause.NewWith("cte", db.Table("users").Where("`name` = ?", "WinterYukky"))).Table("cte").
//	 	Scan(&users)
//
// If you need more advanced WITH clause, you can see With struct.
func NewWith(name string, subquery interface{}, materialized bool, args ...interface{}) With {
	switch v := subquery.(type) {
	case *gorm.DB:
		return With{
			CTEs: []CTE{
				{
					Name:     name,
					Subquery: Subquery{DB: v},
				},
			},
			Materialized: materialized,
		}
	case string:
		return With{
			CTEs: []CTE{
				{
					Name:     name,
					Subquery: clause.Expr{SQL: v, Vars: args},
				},
			},
			Materialized: materialized,
		}
	}

	return With{}
}
