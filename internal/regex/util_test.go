package regex_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeString(t *testing.T) {
	t.Parallel()

	type parseTest struct {
		inputString    string
		expectedOutput string
	}

	parseTests := []parseTest{
		{
			inputString:    "Mission.Impossible 'quoted string' and \"double quoted string\" &&jF $$ q",
			expectedOutput: "mission impossible 'quoted string' and \"double quoted string\" jf q",
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedOutput, regex.NormalizeString(test.inputString))
		})
	}
}

func TestNormalizeSearchString(t *testing.T) {
	t.Parallel()

	type parseTest struct {
		inputString    string
		expectedOutput string
	}

	parseTests := []parseTest{
		{
			inputString:    "Mission.Impossible 'quoted string' and \"double quoted string\" &&jF $$ -weak q",
			expectedOutput: "mission impossible 'quoted string' and \"double quoted string\" jf -weak q",
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedOutput, regex.NormalizeSearchString(test.inputString))
		})
	}
}
