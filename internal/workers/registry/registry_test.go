package registry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/testutils"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRegistry(t *testing.T) {
	t.Parallel()

	test1 := testutils.TestRunner{
		StartupSleep:  10 * time.Millisecond,
		ShutdownSleep: 10 * time.Millisecond,
	}
	test2 := test1
	r, err := registry.NewRegistry(
		zap.NewNop(),
		registry.WithWorker("test1", test1.Runner()),
		registry.WithWorker("test2", test2.Runner()),
	)
	require.NoError(t, err)

	assert.Equal(t, map[string]worker.StateInfo{
		"test1": {
			State: worker.StateIdle,
		},
		"test2": {
			State: worker.StateIdle,
		},
	}, r.WorkersState())

	err = r.Start(t.Context(), "test1", "test2")
	require.NoError(t, err)
	assert.Equal(t, map[string]worker.StateInfo{
		"test1": {
			State: worker.StateRunning,
		},
		"test2": {
			State: worker.StateRunning,
		},
	}, r.WorkersState())

	err = r.Start(t.Context(), "test1")
	require.NoError(t, err)

	<-time.After(10 * time.Millisecond)

	assert.Equal(t, map[string]worker.StateInfo{
		"test1": {
			State: worker.StateRunning,
		},
		"test2": {
			State: worker.StateRunning,
		},
	}, r.WorkersState())

	require.NoError(t, r.Shutdown(t.Context(), "test1"))

	assert.Equal(t, map[string]worker.StateInfo{
		"test1": {
			State: worker.StateIdle,
		},
		"test2": {
			State: worker.StateRunning,
		},
	}, r.WorkersState())

	require.NoError(t, r.Shutdown(t.Context(), "test1", "test2"))
	assert.Equal(t, map[string]worker.StateInfo{
		"test1": {
			State: worker.StateIdle,
		},
		"test2": {
			State: worker.StateIdle,
		},
	}, r.WorkersState())

	test2.StartupErr = assert.AnError

	err = r.Start(t.Context(), "test1", "test2")
	require.Error(t, err)
	require.ErrorIs(t, err, worker.ErrStart)
	require.ErrorIs(t, err, registry.ErrPartial)
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, map[string]worker.StateInfo{
		"test1": {
			State: worker.StateRunning,
		},
		"test2": {
			State: worker.StateError,
			Err:   fmt.Errorf("%w: %w", registry.ErrStart, assert.AnError),
		},
	}, r.WorkersState())

	require.NoError(t, r.Shutdown(t.Context(), "test1"))
	require.Error(t, r.Shutdown(t.Context(), "test2"))
}
