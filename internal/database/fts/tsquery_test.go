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
		{"2 words", "foo Bar", "foo & bar"},
		{"3 words", "foo bar baz", "foo & bar & baz"},
		{"quotes, operators, parens, wildcards", "\"make me a \" . (sandwich | panini) !cheese mayo*",
			"make <-> me <-> a <-> (sandwich | panini) & ! cheese & mayo:*"},
		{"unmatched quotes", "\"make me a sandwich", "make <-> me <-> a <-> sandwich"},
		{"unmatched parens", "\"make me a \" . (sandwich | panini",
			"make <-> me <-> a <-> (sandwich | panini)"},
		{"Ukrainian", "зроби мені бутерброд", "zrobi & meni & buterbrod"},
		{"Chinese", "给我做一个三明治", "Gei <-> Wo <-> Zuo <-> Yi <-> Ge <-> San <-> Ming <-> Zhi"},
		{"Arabic", "اصنع لي شطيرة", "'Sn`' & ly & 'shTyr@'"},
		{"Arabic (quoted)", "\"اصنع لي شطيرة\"", "'Sn`' <-> ly <-> 'shTyr@'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, AppQueryToTsquery(tt.input))
		})
	}
}
