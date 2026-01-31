package gen_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGen(t *testing.T) {
	t.Parallel()

	require.NoError(t, exec.CommandContext(t.Context(), "./fixtures/gen.sh").Run())
}
