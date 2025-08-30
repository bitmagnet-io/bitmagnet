package ecma262

import (
	"regexp"
	"strings"
)

// ECMA262IncompatibilityReason represents why a pattern is incompatible with ECMA 262
type ECMA262IncompatibilityReason string

const (
	ReasonTextBoundaries    ECMA262IncompatibilityReason = "uses \\A or \\z text boundaries"
	ReasonByteMatching      ECMA262IncompatibilityReason = "uses \\C byte matching"
	ReasonPOSIXCharClasses  ECMA262IncompatibilityReason = "uses POSIX character classes [:name:]"
	ReasonUnicodeProperties ECMA262IncompatibilityReason = "uses Unicode property classes \\p{}"
	ReasonInlineCaseFlags   ECMA262IncompatibilityReason = "uses inline case-insensitive groups (?i:)"
	ReasonMissingLowerBound ECMA262IncompatibilityReason = "uses repetition with missing lower bound {,n}"
)

var posixClasses = []string{
	"[:alnum:]", "[:alpha:]", "[:ascii:]", "[:blank:]", "[:cntrl:]",
	"[:digit:]", "[:graph:]", "[:lower:]", "[:print:]", "[:punct:]",
	"[:space:]", "[:upper:]", "[:word:]", "[:xdigit:]",
}

// IsRegexpCompatible checks if a regexp is compatible with ECMA 262 (JavaScript) regex.
func IsRegexpCompatible(re *regexp.Regexp) (bool, []ECMA262IncompatibilityReason) {
	return IsPatternStringCompatible(re.String())
}

// IsPatternStringCompatible checks if a pattern string is compatible with ECMA 262.
func IsPatternStringCompatible(pattern string) (bool, []ECMA262IncompatibilityReason) {
	var reasons []ECMA262IncompatibilityReason

	// Check for text boundaries (\A, \z)
	if strings.Contains(pattern, "\\A") || strings.Contains(pattern, "\\z") {
		reasons = append(reasons, ReasonTextBoundaries)
	}

	// Check for byte matching (\C)
	if strings.Contains(pattern, "\\C") {
		reasons = append(reasons, ReasonByteMatching)
	}

	// Check for POSIX character classes
	for _, class := range posixClasses {
		if strings.Contains(pattern, class) {
			reasons = append(reasons, ReasonPOSIXCharClasses)
			break
		}
	}

	// Check for Unicode property classes (\p{}, \P{})
	if strings.Contains(pattern, "\\p{") || strings.Contains(pattern, "\\P{") {
		reasons = append(reasons, ReasonUnicodeProperties)
	}

	// Check for inline case-insensitive flags
	if strings.Contains(pattern, "(?i:") || strings.Contains(pattern, "(?-i:") {
		reasons = append(reasons, ReasonInlineCaseFlags)
	}

	// Check for repetition with missing lower bound {,n}
	if regexp.MustCompile(`\{,\d+\}`).MatchString(pattern) {
		reasons = append(reasons, ReasonMissingLowerBound)
	}

	return len(reasons) == 0, reasons
}
