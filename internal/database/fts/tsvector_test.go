package fts_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTsvector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		wantTsv fts.Tsvector
		wantStr string
	}{
		{
			input: " 'a':1A bb:2b 'cc ccc':3C  'dD''Dd''':4D e a bb:5 ",
			wantTsv: fts.Tsvector{
				"a": {
					1: 'A',
				},
				"bb": {
					2: 'B',
					5: 'D',
				},
				"cc ccc": {
					3: 'C',
				},
				"dD'Dd'": {
					4: 'D',
				},
				"e": {},
			},
			wantStr: "'a':1A 'bb':2B,5 'cc ccc':3C 'dD''Dd''':4 'e'",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			got, err := fts.ParseTsvector(test.input)

			require.NoError(t, err)
			assert.Equal(t, test.wantTsv, got)
			assert.Equal(t, test.wantStr, got.String())
		})
	}
}
