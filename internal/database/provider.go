package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/exclause"
	"github.com/bitmagnet-io/bitmagnet/internal/database/logger"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type PoolProvider interface {
	Pool() (*pgxpool.Pool, error)
}

type SQLDBProvider interface {
	SQLDB() (*sql.DB, error)
}

type GormDBProvider interface {
	GormDB() (*gorm.DB, error)
}

type DaoProvider interface {
	Dao() (*dao.Query, error)
}

type DaoTransactionProvider interface {
	DaoProvider
	DaoTransaction(func(tx *dao.Query) error) error
}

type Provider interface {
	PoolProvider
	SQLDBProvider
	GormDBProvider
	DaoTransactionProvider
	IsActive() bool
}

type RunnerProvider interface {
	Provider
	runner.Provider
}

type instances struct {
	pool   *pgxpool.Pool
	sqlDB  *sql.DB
	gormDB *gorm.DB
	dao    *dao.Query
}

func New(dsn string, plugins []gorm.Plugin, logger *zap.Logger) RunnerProvider {
	return &provider{
		dsn:     dsn,
		logger:  logger,
		plugins: plugins,
	}
}

type provider struct {
	mtx       sync.RWMutex
	dsn       string
	plugins   []gorm.Plugin
	instances *instances
	err       error
	logger    *zap.Logger
}

func (p *provider) Runner() runner.Runner {
	return func(
		ctx context.Context,
		cancel context.CancelCauseFunc,
	) (shutdowner runner.Shutdowner, err error) {
		shutdowner = runner.NopShutdowner

		var inst instances

		p.mtx.Lock()

		defer func() {
			p.err = err
			if err != nil {
				p.instances = nil

				if inst.pool != nil {
					inst.pool.Close()
				}

				cancel(err)
			} else {
				p.instances = &inst
			}

			p.mtx.Unlock()
		}()

		if p.instances != nil {
			err = fmt.Errorf("%w: %w", Err, runner.ErrAlreadyRunning)

			return
		}

		inst.pool, err = pgxpool.New(ctx, p.dsn)
		if err != nil {
			err = fmt.Errorf("%w: %w", Err, err)

			return
		}

		p.logger.Debug("waiting for database to be ready")

		err = waitForPool(ctx, inst.pool, p.logger)
		if err != nil {
			err = fmt.Errorf("%w: %w", Err, err)

			return
		}

		p.logger.Debug("database is ready")

		inst.sqlDB = stdlib.OpenDBFromPool(inst.pool)

		inst.gormDB, err = gorm.Open(
			postgres.New(postgres.Config{
				Conn: inst.sqlDB,
			}),
			&gorm.Config{
				Logger: logger.New(
					p.logger,
					logger.Config{
						LogLevel:      gormlogger.Info,
						SlowThreshold: time.Second * 30,
					},
				),
				DisableAutomaticPing: true,
			},
		)

		err = inst.gormDB.Use(exclause.New())
		if err != nil {
			err = fmt.Errorf("%w: %w", Err, err)

			return
		}

		for _, plugin := range p.plugins {
			err = inst.gormDB.Use(plugin)
			if err != nil {
				err = fmt.Errorf("%w: plugin failed: %w", Err, err)

				return
			}
		}

		inst.dao = dao.Use(inst.gormDB)

		shutdowner = func(context.Context) error {
			p.mtx.Lock()
			defer p.mtx.Unlock()

			inst.pool.Close()

			p.instances = nil

			return nil
		}

		return
	}
}

func waitForPool(ctx context.Context, pool *pgxpool.Pool, logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	waitFor := time.Duration(0)

	var err error

outer:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()

			break outer
		case <-time.After(waitFor):
			err = pool.Ping(ctx)
			if err == nil {
				break outer
			}

			logger.Warn("ping failed, waiting to try again", zap.Error(err))
		}

		waitFor = time.Second
	}

	if err != nil {
		err = fmt.Errorf("%w: %w", ErrPingFailed, err)
	}

	return err
}

func (p *provider) IsActive() bool {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.instances != nil
}

func (p *provider) Pool() (*pgxpool.Pool, error) {
	inst, err := p.getInstances()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}

	return inst.pool, nil
}

func (p *provider) SQLDB() (*sql.DB, error) {
	inst, err := p.getInstances()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}

	return inst.sqlDB, nil
}

func (p *provider) GormDB() (*gorm.DB, error) {
	inst, err := p.getInstances()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}

	return inst.gormDB, nil
}

func (p *provider) Dao() (*dao.Query, error) {
	inst, err := p.getInstances()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}

	return inst.dao, nil
}

func (p *provider) DaoTransaction(fn func(tx *dao.Query) error) error {
	dao, err := p.Dao()
	if err != nil {
		return err
	}

	return dao.Transaction(fn)
}

func (p *provider) getInstances() (*instances, error) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.err != nil {
		return nil, p.err
	}

	if p.instances == nil {
		return nil, ErrUninitialized
	}

	return p.instances, nil
}
