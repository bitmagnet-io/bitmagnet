package migrator

import (
	"context"
	"database/sql"

	migrationssql "github.com/bitmagnet-io/bitmagnet/migrations"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func New(db *sql.DB, logger *zap.Logger) Migrator {
	initGoose(logger)

	return &migrator{
		db:     db,
		logger: logger,
	}
}

func initGoose(logger *zap.Logger) {
	goose.SetLogger(gooseLogger{logger})
	goose.SetBaseFS(migrationssql.FS)

	err := goose.SetDialect("postgres")
	if err != nil {
		panic(err)
	}
}

type Migrator interface {
	Up(ctx context.Context) error
	UpTo(ctx context.Context, version int64) error
	Down(ctx context.Context) error
	DownTo(ctx context.Context, version int64) error
}

type migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

func (m *migrator) Up(ctx context.Context) error {
	m.logger.Info("checking and applying migrations...")
	return goose.UpContext(ctx, m.db, ".")
}

func (m *migrator) UpTo(ctx context.Context, version int64) error {
	return goose.UpToContext(ctx, m.db, ".", version)
}

func (m *migrator) Down(ctx context.Context) error {
	return goose.DownContext(ctx, m.db, ".")
}

func (m *migrator) DownTo(ctx context.Context, version int64) error {
	return goose.DownToContext(ctx, m.db, ".", version)
}
