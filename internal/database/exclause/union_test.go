package exclause

import (
	"database/sql/driver"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestUnion_Query(t *testing.T) {
	tests := []struct {
		name      string
		operation func(db *gorm.DB) *gorm.DB
		want      string
		wantArgs  []driver.Value
	}{
		{
			name: "When statement is clause.Expr, then should be used as statement",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.Table("general_users").
					Clauses(Union{
						Statements: []clause.Expression{
							clause.Expr{
								SQL:  "ALL ?",
								Vars: []interface{}{db.Table("admin_users")},
							},
						},
					}).Scan(nil)
			},
			want:     "SELECT * FROM `general_users` UNION ALL SELECT * FROM `admin_users`",
			wantArgs: []driver.Value{},
		},
		{
			name: "When statement is exclause.Subquery, then should be used as statement",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.Table("general_users").
					Clauses(Union{
						Statements: []clause.Expression{
							Subquery{
								DB: db.Table("admin_users"),
							},
						},
					}).Scan(nil)
			},
			want:     "SELECT * FROM `general_users` UNION SELECT * FROM `admin_users`",
			wantArgs: []driver.Value{},
		},
		{
			name: "When has multiple UNION, then should be used all UNION clause",
			operation: func(db *gorm.DB) *gorm.DB {
				return db.Table("general_users").
					Clauses(Union{
						Statements: []clause.Expression{
							Subquery{
								DB: db.Table("admin_users"),
							},
						},
					}).
					Clauses(Union{
						Statements: []clause.Expression{
							Subquery{
								DB: db.Table("guest_users"),
							},
						},
					}).Scan(nil)
			},
			want:     "SELECT * FROM `general_users` UNION SELECT * FROM `admin_users` UNION SELECT * FROM `guest_users`",
			wantArgs: []driver.Value{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db, _ := gorm.Open(mysql.New(mysql.Config{
				Conn:                      mockDB,
				SkipInitializeWithVersion: true,
			}))
			if err := db.Use(New()); err != nil {
				t.Fatalf("an error '%s' was not expected when using the database plugin", err)
			}
			mock.ExpectQuery(regexp.QuoteMeta(tt.want)).WithArgs(tt.wantArgs...).WillReturnRows(sqlmock.NewRows([]string{}))
			if tt.operation != nil {
				db = tt.operation(db)
			}
			if db.Error != nil {
				t.Error(db.Error.Error())
			}
		})
	}
}

func TestNewUnion(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))
	db = db.Table("users")
	type args struct {
		subquery interface{}
		args     []interface{}
	}
	tests := []struct {
		name string
		args args
		want Union
	}{
		{
			name: "When subquery is *gorm.DB, then statement is exclause.Subquery",
			args: args{
				subquery: db,
			},
			want: Union{
				Statements: []clause.Expression{
					Subquery{
						DB: db,
					},
				},
			},
		},
		{
			name: "When subquery is string, then statement is clause.Expr",
			args: args{
				subquery: "ALL ?",
				args:     []interface{}{db.Table("users")},
			},
			want: Union{
				Statements: []clause.Expression{
					clause.Expr{
						SQL:  "ALL ?",
						Vars: []interface{}{db.Table("users")},
					},
				},
			},
		},
		{
			name: "When subquery is else, then statement is empty Union",
			args: args{
				subquery: 0,
			},
			want: Union{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnion(tt.args.subquery, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}
