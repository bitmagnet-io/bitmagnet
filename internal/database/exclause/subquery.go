package exclause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Subquery is subquery statement
type Subquery struct {
	DB *gorm.DB
}

// Build build subquery
func (subquery Subquery) Build(builder clause.Builder) {
	builder.AddVar(builder, subquery.DB)
}
