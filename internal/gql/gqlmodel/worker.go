package gqlmodel

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
)

type WorkerQuery struct {
	Registry *registry.Registry
}

func (q *WorkerQuery) ListAll(_ context.Context) (gen.WorkerListAllQueryResult, error) {
	return workersListAll(q.Registry)
}

type WorkerMutation struct {
	Context  env.Context
	Registry *registry.Registry
}

func (m *WorkerMutation) Start(_ context.Context, keys []string) (gen.WorkerListAllQueryResult, error) {
	err := m.Registry.Start(m.Context, keys...)
	if err != nil {
		return gen.WorkerListAllQueryResult{}, err
	}

	return workersListAll(m.Registry)
}

func (m *WorkerMutation) Shutdown(ctx context.Context, keys []string) (gen.WorkerListAllQueryResult, error) {
	if slices.Contains(keys, http_server.Ref.String()) {
		return gen.WorkerListAllQueryResult{}, fmt.Errorf(
			`"%s" worker cannot be shutdown via the API`,
			http_server.Ref.String(),
		)
	}

	state := m.Registry.WorkersState()
	dependentKeys := slice.Filter(keys, func(key string) bool {
		if info, ok := state[key]; ok {
			if _, ok := info.RequiredBy[http_server.Ref.String()]; ok {
				return true
			}
		}

		return false
	})

	if len(dependentKeys) > 0 {
		return gen.WorkerListAllQueryResult{}, fmt.Errorf(
			`cannot shutdown workers because they are required by the "%s" worker: "%s"`,
			http_server.Ref.String(),
			strings.Join(keys, `", "`),
		)
	}

	err := m.Registry.Shutdown(ctx, keys...)
	if err != nil {
		return gen.WorkerListAllQueryResult{}, err
	}

	return workersListAll(m.Registry)
}

func (m *WorkerMutation) Restart(_ context.Context, keys []string) (gen.WorkerListAllQueryResult, error) {
	// Must be done in a goroutine to prevent deadlock:
	go func() {
		_ = m.Registry.Restart(m.Context, keys...)
	}()

	// Hopefully give workers time to enter shutdown state:
	<-time.After(time.Millisecond * 100)

	return workersListAll(m.Registry)
}

func workersListAll(registry *registry.Registry) (gen.WorkerListAllQueryResult, error) {
	stateMap := registry.WorkersState()
	workers := make([]gen.Worker, 0, len(stateMap))

	for key, state := range stateMap {
		var err *string
		if state.Err != nil {
			*err = state.Err.Error()
		}

		workers = append(workers, gen.Worker{
			Key:        key,
			State:      state.State,
			Error:      err,
			RequiredBy: state.RequiredBy.Slice(),
			DependsOn:  state.DependsOn.Slice(),
		})
	}

	slices.SortFunc(workers, func(a, b gen.Worker) int {
		return cmp.Compare(a.Key, b.Key)
	})

	return gen.WorkerListAllQueryResult{
		Workers: workers,
	}, nil
}
