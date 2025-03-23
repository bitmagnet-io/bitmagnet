package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDateRangeFromString(t *testing.T) {
	type parseTest struct {
		inputString   string
		expectedStart Date
		expectedEnd   Date
	}

	var parseTests = []parseTest{
		{
			inputString:   "2020-01-01",
			expectedStart: Date{Year: 2020, Month: 1, Day: 1},
			expectedEnd:   Date{Year: 2020, Month: 1, Day: 1},
		},
		{
			inputString:   "2020-01-02 to 2022",
			expectedStart: Date{Year: 2020, Month: 1, Day: 2},
			expectedEnd:   Date{Year: 2022, Month: 12, Day: 31},
		},
		{
			inputString:   "2020-02 TO 2022-03-04",
			expectedStart: Date{Year: 2020, Month: 2, Day: 1},
			expectedEnd:   Date{Year: 2022, Month: 3, Day: 4},
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			actualOutput, err := NewDateRangeFromString(test.inputString)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedStart, actualOutput.Start())
			assert.Equal(t, test.expectedEnd, actualOutput.End())
		})
	}
}
