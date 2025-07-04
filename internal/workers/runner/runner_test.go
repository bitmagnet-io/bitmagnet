package runner_test

import (
	"context"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/testutils"
	"github.com/stretchr/testify/require"
)

func TestSingleWorker(t *testing.T) {
	t.Parallel()

	for _, test := range []testutils.TestRunner{
		testutils.RunnerSimple,
		testutils.RunnerNilShutdowner,
	} {
		func(test testutils.TestRunner) {
			t.Run(test.Name, func(t *testing.T) {
				t.Parallel()

				runner := test.Runner()

				ctx, cancel := context.WithCancelCause(t.Context())
				t.Cleanup(func() { cancel(nil) })

				shutdown, err := runner(ctx, cancel)
				if test.StartupErr == nil {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, test.StartupErr)
				}

				err = shutdown.Call(ctx)
				if test.ShutdownErr == nil {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, test.ShutdownErr)
				}
			})
		}(test)
	}
}
