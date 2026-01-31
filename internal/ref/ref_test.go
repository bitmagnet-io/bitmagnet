package ref_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	r, err := ref.Parse("core.logging.console")
	require.NoError(t, err)

	assert.Equal(t, "core.logging.console", r.String())
}

func TestSub(t *testing.T) {
	t.Parallel()

	root, err := ref.New("root")
	require.NoError(t, err)

	rootCopy := root

	sub, err := rootCopy.Sub("sub")
	require.NoError(t, err)

	assert.Equal(t, "sub", sub.String())

	_, err = root.Sub("sub")
	require.ErrorIs(t, err, ref.ErrNameAlreadyExists)

	_, err = rootCopy.Sub("sub")
	require.ErrorIs(t, err, ref.ErrNameAlreadyExists)
}
