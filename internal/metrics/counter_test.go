package metrics_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serviceTest = "test_service"
	labelA      = metrics.Label("test_a")
	labelB      = metrics.Label("test_b")
)

var (
	refRoot = ref.MustNew("root")

	labelValueAFoo = labelA.Value("foo")
	labelValueABar = labelA.Value("bar")
	labelValueBFoo = labelB.Value("foo")
	labelValueBBar = labelB.Value("bar")
)

func TestMetricsCounter(t *testing.T) {
	t.Parallel()

	registry, err := metrics.NewRegistry(
		serviceTest,
	)
	require.NoError(t, err)

	defer registry.Close()

	component, err := registry.NewComponent(refRoot)
	require.NoError(t, err)

	counter, err := component.NewCounter("test")
	require.NoError(t, err)

	counter.Incr(labelValueAFoo, labelValueBBar)
	counter.IncrN(2, labelValueABar, labelValueBFoo)

	assert.Equal(t, 3, counter.Total())
	assert.Equal(t, 2, counter.Total(labelValueABar))
	assert.Equal(t, map[string]int{
		"bar": 2,
		"foo": 1,
	}, counter.TotalsByLabel(labelA))
	assert.Equal(t, map[string]int{
		"bar": 1,
		"foo": 2,
	}, counter.TotalsByLabel(labelB))
}
