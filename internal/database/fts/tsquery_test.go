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
		{"operators", "\"make me a\" . (sandwich | panini) !cheese", "make <-> me <-> a <-> (sandwich | panini) & ! cheese"},
		{"Ukrainian", "зроби мені бутерброд", "zrobi & meni & buterbrod"},
		{"Chinese", "给我做一个三明治", "Gei <-> Wo <-> Zuo <-> Yi <-> Ge <-> San <-> Ming <-> Zhi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, AppQueryToTsquery(tt.input))
		})
	}
}
