package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExtensionFromPath(t *testing.T) {
	t.Parallel()

	type parseTest struct {
		inputString string
		expectedExt string
	}

	var parseTests = []parseTest{
		{
			inputString: "short film.mkv",
			expectedExt: "mkv",
		},
		{
			inputString: "short film.mkv.mk2",
			expectedExt: "mk2",
		},
		{
			inputString: "short film....mkv",
			expectedExt: "mkv",
		},
		{
			inputString: ".mkv",
			expectedExt: "",
		},
		{
			inputString: "short film.",
			expectedExt: "",
		},
	}

	for _, test := range parseTests {
		t.Run(test.inputString, func(t *testing.T) {
			t.Parallel()

			actualOutput := FileExtensionFromPath(test.inputString)
			assert.Equal(t, test.expectedExt, actualOutput.String)
		})
	}
}
