package regex

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/hedhyw/rex/pkg/dialect/base"
	"github.com/hedhyw/rex/pkg/rex"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func AnyWordChar() base.ClassToken {
	return rex.Common.Class(rex.Chars.Unicode(unicode.L), rex.Chars.Digits())
}

func AnyNonWordChar() base.ClassToken {
	return rex.Common.NotClass(rex.Chars.Unicode(unicode.L), rex.Chars.Digits())
}

func OpeningPunctuationToken() base.ClassToken {
	return rex.Chars.Runes("('\"")
}

func ClosingPunctuationToken() base.ClassToken {
	return rex.Chars.Runes(",;:?!-)'\"")
}

func MidWordPunctuationToken() base.ClassToken {
	return rex.Chars.Runes("'-")
}

func TrimmedWordToken() base.GroupToken {
	return rex.Group.Composite(
		rex.Group.NonCaptured(rex.Chars.Upper(), rex.Chars.Single('.')).Repeat().EqualOrMoreThan(2),
		rex.Group.NonCaptured(
			AnyWordChar().Repeat().OneOrMore(),
			rex.Group.NonCaptured(
				MidWordPunctuationToken().Repeat().OneOrMore(), AnyWordChar().Repeat().OneOrMore(),
			).Repeat().ZeroOrMore(),
		),
	).NonCaptured()
}

func WordToken() base.GroupToken {
	return rex.Group.Define(
		OpeningPunctuationToken().Repeat().ZeroOrMore(),
		TrimmedWordToken(),
		ClosingPunctuationToken().Repeat().ZeroOrMore(),
	).NonCaptured()
}

var wordTokenRegex = rex.New(WordToken()).MustCompile()

func WordTokenRegex() *regexp.Regexp {
	return wordTokenRegex
}

func NormalizeString(input string) string {
	input = strings.ToLower(input)
	input, _, _ = transform.String(transform.Chain(norm.NFD, norm.NFC), input)

	var tokens []string

	for _, match := range WordTokenRegex().FindAllStringSubmatch(input, -1) {
		if len(match) >= 1 && len(match[0]) >= 1 {
			tokens = append(tokens, match[0])
		}
	}

	return strings.Join(tokens, " ")
}

func QuotedStringToken(quoteCharToken base.ClassToken) base.GroupToken {
	return rex.Group.Define(
		quoteCharToken,
		rex.Group.Composite(
			rex.Group.Define(rex.Chars.Runes("\\"), quoteCharToken).NonCaptured(),
			rex.Common.NotClass(quoteCharToken),
		).Repeat().ZeroOrMore(),
		quoteCharToken,
	).NonCaptured()
}

var searchTokenRegex = rex.New(
	rex.Group.Composite(
		rex.Group.Define(QuotedStringToken(rex.Chars.Runes("'"))),
		rex.Group.Define(QuotedStringToken(rex.Chars.Runes("\""))),
		rex.Group.Define(rex.Chars.Runes(`-`).Repeat().ZeroOrMore(), WordToken()),
	).NonCaptured().Repeat().ZeroOrMore(),
).MustCompile()

func SearchStringToNormalizedTokens(input string) []string {
	var tokens []string

	matches := searchTokenRegex.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		switch {
		case len(match[1]) >= 1:
			tokens = append(tokens, match[1])
		case len(match[3]) >= 1:
			tokens = append(tokens, match[3])
		case len(match[5]) >= 1:
			tokens = append(tokens, strings.ToLower(match[5]))
		}
	}

	return tokens
}

func NormalizeSearchString(input string) string {
	return strings.Join(SearchStringToNormalizedTokens(input), " ")
}
