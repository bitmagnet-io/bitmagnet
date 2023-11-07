package music

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
	"strconv"
)

type CoreInfo struct {
	Title  string
	Album  string
	Artist string
	Year   model.Year
}

var titleTokens = []dialect.Token{
	rex.Group.Define(
		rex.Group.Composite(
			rex.Group.NonCaptured(
				regex.AnyWordChar().Repeat().OneOrMore(),
				rex.Group.NonCaptured(
					rex.Chars.Single('-'), regex.AnyWordChar().Repeat().OneOrMore(),
				).Repeat().ZeroOrMore(),
			),
			regex.AnyNonWordChar().Repeat().OneOrMore(),
		).NonCaptured().Repeat().OneOrMore(),
		rex.Group.Composite(
			regex.AnyNonWordChar().Repeat().OneOrMore(),
			rex.Chars.End(),
		).NonCaptured(),
	),
}

var titleRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(titleTokens...),
).MustCompile()

var yearTokens = []dialect.Token{
	rex.Group.NonCaptured(rex.Common.NotClass(rex.Chars.WordCharacter()).Repeat().ZeroOrMore()),
	rex.Group.Define(
		rex.Group.Composite(
			rex.Common.Text("18"), rex.Common.Text("19"), rex.Common.Text("20"),
		).NonCaptured(),
		rex.Chars.Digits().Repeat().Exactly(2),
	),
	rex.Group.Composite(
		rex.Common.NotClass(rex.Chars.WordCharacter()),
		rex.Chars.End(),
	).NonCaptured(),
}

var titleYearRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(rex.Group.NonCaptured(titleTokens...), rex.Group.NonCaptured(yearTokens...)),
).MustCompile()

var artistDiscographyRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.Define(
		rex.Common.NotClass(rex.Chars.Single('-')).Repeat().OneOrMore()),
	rex.Chars.Whitespace().Repeat().ZeroOrMore(),
	rex.Chars.Single('-').Repeat().ZeroOrMore(),
	rex.Chars.Whitespace().Repeat().ZeroOrMore(),
	rex.Common.Text("discography"),
).MustCompile()

var separatorToken = rex.Chars.Runes(" ._")

var titlePartRegex = rex.New(
	separatorToken.Repeat().ZeroOrOne(),
	rex.Group.Define(regex.WordToken()),
	separatorToken.Repeat().ZeroOrOne(),
).MustCompile()

var trimTitleRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(
		rex.Chars.Single('['),
		rex.Common.NotClass(rex.Chars.Single(']')).Repeat().OneOrMore(),
		rex.Chars.Single(']'),
	).Repeat().ZeroOrOne(),
	regex.AnyNonWordChar().Repeat().ZeroOrMore(),
	rex.Group.Define(
		regex.WordToken(),
		rex.Group.NonCaptured(
			rex.Chars.Any(),
			regex.WordToken(),
		).Repeat().ZeroOrMore(),
	),
	regex.AnyNonWordChar().Repeat().ZeroOrMore(),
	rex.Chars.End(),
).MustCompile()

func cleanTitle(title string) string {
	title = titlePartRegex.ReplaceAllStringFunc(title, func(s string) string {
		partMatch := titlePartRegex.FindStringSubmatch(s)
		if partMatch == nil {
			return ""
		}
		return partMatch[1] + " "
	})
	title = trimTitleRegex.ReplaceAllString(title, "$1")
	return title
}

func parseTitleYear(input string) (CoreInfo, string, error) {
	if match := titleYearRegex.FindStringSubmatch(input); match != nil {
		yearMatch, _ := strconv.ParseUint(match[2], 10, 16)
		title := cleanTitle(match[1])
		if title != "" {
			return CoreInfo{
				Title: title,
				Year:  model.Year(yearMatch),
			}, input[len(match[0]):], nil
		}
	}
	return CoreInfo{}, "", resolver.ErrNoMatch
}

func parseTitle(input string) (CoreInfo, string, error) {
	if match := titleRegex.FindStringSubmatch(input); match != nil {
		title := cleanTitle(match[1])
		if title != "" {
			return CoreInfo{
				Title: title,
			}, input[len(match[0]):], nil
		}
	}
	return CoreInfo{}, "", resolver.ErrNoMatch
}

func FindArtistDiscography(input string) (string, error) {
	if match := artistDiscographyRegex.FindStringSubmatch(input); match != nil {
		return match[1], nil
	} else {
		return "", resolver.ErrNoMatch
	}
}
