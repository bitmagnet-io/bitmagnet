package ecma262

import "regexp"

// Compile compiles a regular expression that is interoperable between RE2 and ECMA 262.
func Compile(pattern string) (*regexp.Regexp, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return re, RegexpCompatibilityError(re)
}

// MustCompile compiles a regular expression that is interoperable between RE2 and ECMA 262, and panics on error.
func MustCompile(pattern string) *regexp.Regexp {
	re, err := Compile(pattern)
	if err != nil {
		panic(err)
	}

	return re
}
