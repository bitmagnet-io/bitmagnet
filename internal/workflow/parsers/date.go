package parsers

import (
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
	"regexp"
)

var DateRegex *regexp.Regexp

func init() {
	var groups []dialect.Token
	for _, sep := range []rune{'-', '.', '/'} {
		groups = append(groups, rex.Group.Define(
			rex.Chars.Digits().Repeat().Exactly(4),
			rex.Chars.Single(sep),
			rex.Chars.Digits().Repeat().Between(1, 2),
			rex.Chars.Single(sep),
			rex.Chars.Digits().Repeat().Between(1, 2),
		).NonCaptured())
		groups = append(groups, rex.Group.Define(
			rex.Chars.Digits().Repeat().Between(1, 2),
			rex.Chars.Single(sep),
			rex.Chars.Digits().Repeat().Between(1, 2),
			rex.Chars.Single(sep),
			rex.Chars.Digits().Repeat().Exactly(4),
		).NonCaptured())
	}
	DateRegex = rex.New(
		rex.Group.Composite(groups...),
	).MustCompile()
}
