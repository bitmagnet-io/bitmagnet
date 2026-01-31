package metrics_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsGauge(t *testing.T) {
	t.Parallel()

	registry, err := metrics.NewRegistry("test")
	require.NoError(t, err)

	defer registry.Close()

	component, err := registry.NewComponent(refRoot)
	require.NoError(t, err)

	gauge, err := component.NewGauge("test")
	require.NoError(t, err)

	// The gauge should return values for the exact combination of labels that were set

	gauge.Set(1, labelValueAFoo)
	gauge.Set(3, labelValueABar)
	gauge.Set(2)

	assert.Equal(t, 2, gauge.Value())
	assert.Equal(t, 1, gauge.Value(labelValueAFoo))
	assert.Equal(t, 0, gauge.Value(labelValueAFoo, labelValueBFoo))
	assert.Equal(t, 3, gauge.Value(labelValueABar))
	assert.Equal(t, 0, gauge.Value(labelValueBFoo))
	assert.Equal(t, 0, gauge.Value(labelValueBBar))

	gauge.Set(3, labelValueAFoo, labelValueBBar)

	assert.Equal(t, map[string]int{"bar": 3, "foo": 1}, gauge.ValuesByLabel(labelA))

	assert.Equal(t, 2, gauge.Value())
	assert.Equal(t, 1, gauge.Value(labelValueAFoo))
	assert.Equal(t, 3, gauge.Value(labelValueAFoo, labelValueBBar))

	// the order of the labels shouldn't matter

	assert.Equal(t, 0, gauge.Value(labelValueAFoo, labelValueBFoo))
}
