package plugin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/proto/host/configurator"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/tetratelabs/wazero"
	wasi_snapshot_preview1 "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type instanceBuilder struct {
	env           env.Env
	manifest      Manifest
	data          []byte
	runtimeConfig wazero.RuntimeConfig
	moduleConfig  wazero.ModuleConfig
	instantiators []instantiator
}

type instantiator func(ctx context.Context, runtime wazero.Runtime) error

func (p *Plugin) NewInstance(
	env env.Env,
	resolvedConfig resolver.Resolved,
) (Instance, error) {
	builder := &instanceBuilder{
		env:      env,
		manifest: p.manifest,
		data:     p.data,
		runtimeConfig: wazero.NewRuntimeConfig().
			WithCloseOnContextDone(true).
			WithCompilationCache(p.compilationCache),
		moduleConfig: wazero.NewModuleConfig().
			WithStartFunctions("_initialize").
			WithStdout(env).
			WithStderr(env.Stderr()).
			WithSysWalltime().
			WithSysNanotime().
			WithSysNanosleep(),
		instantiators: []instantiator{
			func(ctx context.Context, runtime wazero.Runtime) error {
				_, err := wasi_snapshot_preview1.Instantiate(ctx, runtime)
				return err
			},
		},
	}

	if len(p.manifest.Config) > 0 {
		cfg := make(map[string]any)

		for _, param := range p.configParams {
			if value, ok := resolvedConfig.Param(param.Ref); ok {
				cfg[param.Name()] = value.Value()
			}
		}

		jsonCfg, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}

		builder.instantiators = append(
			builder.instantiators,
			func(ctx context.Context, runtime wazero.Runtime) error {
				return configurator.Instantiate(ctx, runtime, configuratorImpl{
					jsonConfig: string(jsonCfg),
				})
			},
		)
	}

	if p.manifest.Permissions.FS != nil {
		p.manifest.Permissions.FS.build(builder)
	}

	if p.manifest.Permissions.HTTP != nil {
		p.manifest.Permissions.HTTP.build(builder)
	}

	mtx := make(chan struct{}, 1)

	pool := pool.NewObjectPool(
		env,
		pool.NewPooledObjectFactory(
			func(ctx context.Context) (any, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case mtx <- struct{}{}:
				}

				defer func() { <-mtx }()

				return builder.newModule(ctx)
			},
			func(ctx context.Context, obj *pool.PooledObject) error {
				return obj.Object.(*module).Close(ctx)
			},
			func(_ context.Context, obj *pool.PooledObject) bool {
				return !obj.Object.(*module).IsClosed()
			},
			nil,
			nil,
		),
		&pool.ObjectPoolConfig{
			LIFO:                     true,
			MaxTotal:                 10,
			MaxIdle:                  10,
			MinIdle:                  1,
			TestOnBorrow:             true,
			TestWhileIdle:            true,
			BlockWhenExhausted:       true,
			SoftMinEvictableIdleTime: time.Minute,
			TimeBetweenEvictionRuns:  time.Minute,
			EvictionContext:          env,
		},
	)

	pool.PreparePool(env)

	return &instance{
		manifest: p.manifest,
		pool:     pool,
	}, nil
}
