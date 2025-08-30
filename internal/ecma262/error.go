package ecma262

import (
	"errors"
	"fmt"
	"regexp"
)

var Err = errors.New("regular expression is incompatible with ECMA 262")

type ErrReasons []ECMA262IncompatibilityReason

func (rs ErrReasons) Error() string {
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

// PatternStringCompatibilityError returns an error if the provided regular expression pattern string is not compatible with ECMA 262
func PatternStringCompatibilityError(str string) error {
	compatible, reasons := IsPatternStringCompatible(str)

	if compatible {
		return nil
	} else {
		return fmt.Errorf("%w: %s: %w", Err, str, ErrReasons(reasons))
	}
}
