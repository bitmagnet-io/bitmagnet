package ktable

import (
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_bucketIndex(t *testing.T) {
	tests := []struct {
		origin, target string
		index          int
	}{
		{
			"ffbfba8945192d408d3dcc52ba24903a00000000",
			"ffbfba8945192d408d3dcc52ba24903a00000001",
			159,
		},
		{
			"ffcfba8945192d408d3dcc52ba24903a00000000",
			"ffbfba8945192d408d3dcc52ba24903a00000001",
			9,
		},
	}
	for _, test := range tests {
		origin := protocol.MustParseID(test.origin)
		target := protocol.MustParseID(test.target)
		index := test.index
		assert.Equal(t, index, bucketIndex(target, origin))
	}
}
