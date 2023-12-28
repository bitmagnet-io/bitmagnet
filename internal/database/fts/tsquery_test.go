package fts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppQueryToTsquery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"1 word", "foo", "foo"},
		{"2 words", "foo bar", "foo & bar"},
		{"3 words", "foo bar baz", "foo & bar & baz"},
		{"quoted", `"foo bar"`, "foo <-> bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, AppQueryToTsquery(tt.input))
		})
	}
}
