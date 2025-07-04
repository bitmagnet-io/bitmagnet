package concat_test

import (
	"context"
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/concat"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/testutils"
	"github.com/stretchr/testify/require"
)

type concatTest struct {
	name            string
	runners         []testutils.TestRunner
	wantStartupErr  []error
	wantShutdownErr []error
	minRuntime      time.Duration
}

func (t *concatTest) runner() runner.Runner {
	runners := make([]runner.Runner, 0, len(t.runners))
	for _, w := range t.runners {
		runners = append(runners, w.Runner())
	}

	return concat.Runners(runners...)
}

func TestConcat(t *testing.T) {
	t.Parallel()

	for _, test := range []concatTest{
		{
			name: "empty",
		},
		{
			name: "single",
			runners: []testutils.TestRunner{
				testutils.RunnerSimple,
				testutils.RunnerNilShutdowner,
			},
			minRuntime: time.Millisecond * 80,
		},
		{
			name: "double",
			runners: []testutils.TestRunner{
				testutils.RunnerSimple,
				testutils.RunnerNilShutdowner,
			},
		},
		{
			name: "partial failure",
			runners: []testutils.TestRunner{
				testutils.RunnerStartupFailure,
				testutils.RunnerSimple,
				testutils.RunnerNilShutdowner,
			},
			wantStartupErr: []error{
				testutils.RunnerStartupFailure.StartupErr,
				concat.ErrPartial,
			},
		},
	} {
		func(test concatTest) {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				run := test.runner()

				ctx, cancel := context.WithCancelCause(t.Context())
				t.Cleanup(func() { cancel(nil) })

				shutdown, err := run(ctx, cancel)
				if test.wantStartupErr == nil {
					require.NoError(t, err)
				} else {
					for _, wantErr := range test.wantStartupErr {
						require.ErrorIs(t, err, wantErr)
					}

					require.ErrorIs(t, err, concat.Err)
				}

				<-time.After(test.minRuntime)

				require.NoError(t, ctx.Err())

				err = shutdown.Call(ctx)
				if test.wantShutdownErr == nil {
					require.NoError(t, err)
				} else {
					for _, wantErr := range test.wantShutdownErr {
						require.ErrorIs(t, err, wantErr)
					}

					require.ErrorIs(t, err, concat.Err)
					require.ErrorIs(t, err, concat.ErrShutdown)
				}
			})
		}(test)
	}
}
