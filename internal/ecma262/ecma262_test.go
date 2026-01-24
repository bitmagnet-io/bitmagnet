package ecma262_test

import (
	"regexp"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/ecma262"
)

func TestIsCompatible(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		pattern          string
		expectCompatible bool
		expectReasons    []ecma262.IncompatibilityReason
	}{
		{
			name:             "Compatible pattern",
			pattern:          `^[a-zA-Z0-9]+$`,
			expectCompatible: true,
			expectReasons:    nil,
		},
		{
			name:             "Text boundaries",
			pattern:          `\Ahello\z`,
			expectCompatible: false,
			expectReasons:    []ecma262.IncompatibilityReason{ecma262.ReasonTextBoundaries},
		},
		{
			name:             "POSIX character classes",
			pattern:          `[[:alpha:]]+`,
			expectCompatible: false,
			expectReasons:    []ecma262.IncompatibilityReason{ecma262.ReasonPOSIXCharClasses},
		},
		{
			name:             "Unicode property classes",
			pattern:          `\p{L}+`,
			expectCompatible: false,
			expectReasons:    []ecma262.IncompatibilityReason{ecma262.ReasonUnicodeProperties},
		},
		{
			name:             "Inline case flags",
			pattern:          `(?i:hello)`,
			expectCompatible: false,
			expectReasons:    []ecma262.IncompatibilityReason{ecma262.ReasonInlineCaseFlags},
		},
		{
			name:             "Missing lower bound",
			pattern:          `a{,5}`,
			expectCompatible: false,
			expectReasons:    []ecma262.IncompatibilityReason{ecma262.ReasonMissingLowerBound},
		},
		{
			name:             "Multiple incompatibilities",
			pattern:          `\Ahello[[:alpha:]]+\z`,
			expectCompatible: false,
			expectReasons: []ecma262.IncompatibilityReason{
				ecma262.ReasonTextBoundaries,
				ecma262.ReasonPOSIXCharClasses,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			re := regexp.MustCompile(tt.pattern)
			compatible, reasons := ecma262.IsRegexpCompatible(re)

			if compatible != tt.expectCompatible {
				t.Errorf("Expected compatibility %v, got %v", tt.expectCompatible, compatible)
			}

			if len(reasons) != len(tt.expectReasons) {
				t.Errorf(
					"Expected %d reasons, got %d: %v",
					len(tt.expectReasons),
					len(reasons),
					reasons,
				)

				return
			}

			for i, expectedReason := range tt.expectReasons {
				if i >= len(reasons) || reasons[i] != expectedReason {
					t.Errorf(
						"Expected reason %q at index %d, got %q",
						expectedReason,
						i,
						reasons[i],
					)
				}
			}
		})
	}
}
