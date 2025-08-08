package ref_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	t.Parallel()

	r, err := ref.Parse("core.logging.console")
	require.NoError(t, err)

	assert.Equal(t, "core.logging.console", r.String())
}
