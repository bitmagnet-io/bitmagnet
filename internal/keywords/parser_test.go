package keywords

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   []string
		match   []string
		nomatch []string
		err     error
	}{
		{
			input:   []string{"foo", "bar*"},
			match:   []string{"foo", "bar", "barfoo", "x foo.x", "foo/bar"},
			nomatch: []string{"bat", "foobar"},
		},
		{
			input:   []string{"foo?"},
			match:   []string{"foo", "fo"},
			nomatch: []string{"f", "fooo"},
		},
		{
			input: []string{"foo(bar|bat)?", "qux"},
			match: []string{"foo", "foobar", "foobat", "qux"},
		},
		{
			input: []string{"foo(bar"},
			err:   ErrUnexpectedEOF,
		},
		{
			input: []string{"foo\\(bar"},
			match: []string{"foo(bar"},
		},
		{
			input: []string{"(audio)?books?", "(auto)?biograph(y|ies)"},
			match: []string{"book", "books", "audiobook", "audiobooks", "biography", "Autobiographies"},
		},
	}
	for _, test := range tests {
		t.Run(strings.Join(test.input, ", "), func(t *testing.T) {
			t.Parallel()

			r, err := NewRegexFromKeywords(test.input...)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				return
			}

			require.NoError(t, err)
			t.Logf("regex: %s", r.String())

			for _, m := range test.match {
				assert.True(t, r.MatchString(m))
			}

			for _, nm := range test.nomatch {
				assert.False(t, r.MatchString(nm))
			}
		})
	}
}
