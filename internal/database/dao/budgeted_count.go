package dao

import (
	"gorm.io/gorm"
)

func ToSQL(db *gorm.DB) string {
	return db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&[]interface{}{})
	})
}

type BudgetedCountResult struct {
	Count          int64
	Cost           float64
	BudgetExceeded bool
}

func BudgetedCount(db *gorm.DB, budget float64) (BudgetedCountResult, error) {
	row := db.Raw("SELECT count, cost, budget_exceeded from budgeted_count(?, ?)", ToSQL(db), budget).Row()
	result := BudgetedCountResult{}
	err := row.Scan(&result.Count, &result.Cost, &result.BudgetExceeded)
	return result, err
}
