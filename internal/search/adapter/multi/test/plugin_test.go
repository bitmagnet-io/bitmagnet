package test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/search/adapter/multi"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/plugin/test"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestAdapter(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t, zaptest.Level(zapcore.InfoLevel))
	undo := zap.ReplaceGlobals(logger)

	t.Cleanup(func() {
		undo()
	})
	test.BuildTestPlugins(t)

	opts := make([]plugin.ProviderOption, 0, 2)
	for _, p := range []string{"fake_search", "fake_search2"} {
		opts = append(opts, plugin.LoadPlugin(test.PluginTestDataDir(1)+"/"+p, ""))
	}
	bundle, err := plugin.NewProvider(opts...)
	require.NoError(t, err)

	penv := env.NewDefault()

	rplugins, err := registry.New(bundle).Resolve(penv)
	require.NoError(t, err)
	assert.NotNil(t, rplugins)

	type searchParams struct {
		fx.In
		Search []multi.Index `group:"search_adapters"`
	}

	testSearch := func(p searchParams) {
		for _, splugin := range p.Search {
			r, err := (splugin.Adapter).(search.TorrentContent).TorrentContent(t.Context(), search.Params{})
			require.NoError(t, err)
			assert.Equal(t, uint(99), r.TotalCount.Uint)
		}
	}

	fxopt := []fx.Option{
		fx.Provide(
			func() *zap.Logger { return logger },
			func() env.Env { return penv },
		),
		fx.Invoke(
			testSearch,
		),
	}
	app := rplugins.Build(fxopt...)
	require.NoError(t, app.Err())
}
