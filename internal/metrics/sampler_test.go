package metrics_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsSampler(t *testing.T) {
	t.Parallel()

	registry, err := metrics.NewRegistry(
		serviceTest,
	)
	require.NoError(t, err)

	defer registry.Close()

	component, err := registry.NewComponent(refRoot)
	require.NoError(t, err)

	sampler, err := component.NewSampler("test")
	require.NoError(t, err)

	sampler.Add(1, labelValueAFoo, labelValueBBar)
	sampler.Add(3, labelValueABar, labelValueBFoo)

	assert.InDelta(t, 2, sampler.Average().Value(), 0)
	assert.Equal(t, map[string]metrics.Average{
		"bar": {
			Count: 1,
			Sum:   3,
		},
		"foo": {
			Count: 1,
			Sum:   1,
		},
	}, sampler.AveragesByLabel(labelA))
}
