package regex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeString(t *testing.T) {

	type parseTest struct {
		inputString    string
		expectedOutput string
	}

	var parseTests = []parseTest{
		{
			inputString:    "Mission.Impossible 'quoted string' and \"double quoted string\" &&jF $$ q",
			expectedOutput: "mission impossible 'quoted string' and \"double quoted string\" jf q",
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			assert.Equal(t, test.expectedOutput, NormalizeString(test.inputString))
		})
	}
}

func TestNormalizeSearchString(t *testing.T) {

	type parseTest struct {
		inputString    string
		expectedOutput string
	}

	var parseTests = []parseTest{
		{
			inputString:    "Mission.Impossible 'quoted string' and \"double quoted string\" &&jF $$ -weak q",
			expectedOutput: "mission impossible 'quoted string' and \"double quoted string\" jf -weak q",
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			assert.Equal(t, test.expectedOutput, NormalizeSearchString(test.inputString))
		})
	}
}
