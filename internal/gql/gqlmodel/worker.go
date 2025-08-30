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
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
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

func (m *WorkerMutation) Start(_ context.Context, refs []ref.Ref) (gen.WorkerListAllQueryResult, error) {
	err := m.Registry.Start(m.Context, refs...)
	if err != nil {
		return gen.WorkerListAllQueryResult{}, err
	}

	return workersListAll(m.Registry)
}

func (m *WorkerMutation) Shutdown(ctx context.Context, refs []ref.Ref) (gen.WorkerListAllQueryResult, error) {
	if slice.Some(refs, func(ref ref.Ref) bool {
		return ref.String() == http_server.Ref.String()
	}) {
		return gen.WorkerListAllQueryResult{}, fmt.Errorf(
			`"%s" worker cannot be shutdown via the API`,
			http_server.Ref.String(),
		)
	}

	state := m.Registry.WorkersState()
	dependentKeys := slice.Filter(refs, func(ref ref.Ref) bool {
		return state.Get(ref).RequiredBy.Has(http_server.Ref)
	})

	if len(dependentKeys) > 0 {
		return gen.WorkerListAllQueryResult{}, fmt.Errorf(
			`cannot shutdown workers because they are required by the "%s" worker: "%s"`,
			http_server.Ref.String(),
			strings.Join(slice.Map(refs, func(ref ref.Ref) string {
				return ref.String()
			}), `", "`),
		)
	}

	err := m.Registry.Shutdown(ctx, refs...)
	if err != nil {
		return gen.WorkerListAllQueryResult{}, err
	}

	return workersListAll(m.Registry)
}

func (m *WorkerMutation) Restart(_ context.Context, refs []ref.Ref) (gen.WorkerListAllQueryResult, error) {
	// Must be done in a goroutine to prevent deadlock:
	go func() {
		_ = m.Registry.Restart(m.Context, refs...)
	}()

	// Hopefully give workers time to enter shutdown state:
	<-time.After(time.Millisecond * 100)

	return workersListAll(m.Registry)
}

func workersListAll(registry *registry.Registry) (gen.WorkerListAllQueryResult, error) {
	stateMap := registry.WorkersState()
	workers := make([]gen.Worker, 0, stateMap.Len())

	for _, state := range stateMap.Values() {
		var err *string
		if state.Err != nil {
			strErr := state.Err.Error()
			err = &strErr
		}

		workers = append(workers, gen.Worker{
			Ref:        state.Ref,
			State:      state.State,
			Error:      err,
			RequiredBy: state.RequiredBy.Refs(),
			DependsOn:  state.DependsOn.Refs(),
		})
	}

	slices.SortFunc(workers, func(a, b gen.Worker) int {
		return cmp.Compare(a.Ref.String(), b.Ref.String())
	})

	return gen.WorkerListAllQueryResult{
		Workers: workers,
	}, nil
}
