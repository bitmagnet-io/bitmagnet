package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	pool "github.com/jolestar/go-commons-pool/v2"
)

func (p *Plugin) NewInstance(
	env env.Env,
	resolvedConfig resolver.Resolved,
) (Instance, error) {
	cfg := make(map[string]any)

	for _, param := range p.configParams {
		if value, ok := resolvedConfig.Param(param.Ref); ok {
			cfg[param.Ref.Name()] = value.Value()
		}
	}

	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	contract, err := p.api.Configure(env, &api.JSONPayload{
		Data: jsonCfg,
	})
	if err != nil {
		return nil, err
	}

	if contract.GetPermissions() != nil && contract.GetPermissions().GetHttp() != nil {
		p.httpEgress.Set(contract.GetPermissions().GetHttp().GetEgress())
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

				apiModule, err := p.newModule(ctx)
				if err != nil {
					return nil, err
				}

				pluginAPI, err := pluginPlugin.LoadModule(ctx, apiModule)
				if err != nil {
					return nil, fmt.Errorf("failed to load plugin lifecycle: %w", err)
				}

				_, err = pluginAPI.Configure(ctx, &api.JSONPayload{
					Data: jsonCfg,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to configure plugin: %w", err)
				}

				_, err = pluginAPI.Instantiate(ctx, &api.Empty{})
				if err != nil {
					return nil, fmt.Errorf("failed to instantiate plugin: %w", err)
				}

				return apiModule, nil
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

	return &instance{
		ref:      p.ref,
		contract: contract,
		pool:     pool,
	}, nil
}
