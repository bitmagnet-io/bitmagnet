package fts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTsvector(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Tsvector
	}{
		{
			name:  "happy path",
			input: "'a':1A 'bb':2b 'cc ccc':3C 'dD\\'Dd\\'':4D",
			want: Tsvector{
				"a": map[TsvectorLabel]struct{}{
					{1, 'A'}: {},
				},
				"bb": map[TsvectorLabel]struct{}{
					{2, 'B'}: {},
				},
				"cc ccc": map[TsvectorLabel]struct{}{
					{3, 'C'}: {},
				},
				"dD'Dd'": map[TsvectorLabel]struct{}{
					{4, 'D'}: {},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseTsvector(test.input)
			if err != nil {
				t.Errorf("ParseTsvector(%q) = %v", test.input, err)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}
}
