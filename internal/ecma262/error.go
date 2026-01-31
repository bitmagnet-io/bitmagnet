package ecma262

import (
	"errors"
	"fmt"
	"regexp"
)

var Err = errors.New("regular expression is incompatible with ECMA 262")

type IncompatibilityError []IncompatibilityReason

func (rs IncompatibilityError) Error() string {
	msg := ""

	for i, r := range rs {
		if i > 0 {
			msg += "; "
		}

		msg += string(r)
	}

	return msg
}

// RegexpCompatibilityError returns an error if the provided regular expression is not compatible with ECMA 262
func RegexpCompatibilityError(re *regexp.Regexp) error {
	return PatternStringCompatibilityError(re.String())
}

// PatternStringCompatibilityError returns an error if the provided regular expression pattern string is not compatible
// with ECMA 262
func PatternStringCompatibilityError(str string) error {
	compatible, reasons := IsPatternStringCompatible(str)

	if compatible {
		return nil
	}

	return fmt.Errorf("%w: %s: %w", Err, str, IncompatibilityError(reasons))
}
