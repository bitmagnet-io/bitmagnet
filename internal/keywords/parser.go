package keywords

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/dialect/base"
	"github.com/hedhyw/rex/pkg/rex"
)

func NewRegexFromKeywords(kws ...string) (*regexp.Regexp, error) {
	tokens, err := NewRexTokensFromKeywords(kws...)
	if err != nil {
		return nil, err
	}

	return rex.New(
		rex.Group.Composite(
			rex.Chars.Begin(),
			regex.AnyNonWordChar().Repeat().OneOrMore(),
		).NonCaptured(),
		rex.Group.Composite(
			tokens...,
		),
		rex.Group.Composite(
			rex.Chars.End(),
			regex.AnyNonWordChar().Repeat().OneOrMore(),
		).NonCaptured(),
	).Compile()
}

func MustNewRegexFromKeywords(kws ...string) *regexp.Regexp {
	r, err := NewRegexFromKeywords(kws...)
	if err != nil {
		panic(err)
	}

	return r
}

func NewRexTokensFromKeywords(kws ...string) ([]dialect.Token, error) {
	if len(kws) == 0 {
		return nil, errors.New("no keywords provided")
	}

	tokens := make([]dialect.Token, 0, len(kws))

	usedKeywords := make(map[string]struct{})
	for _, kw := range kws {
		if _, ok := usedKeywords[kw]; ok {
			continue
		}

		usedKeywords[kw] = struct{}{}
		l := keywordsLexer{Lexer: lexer.NewLexer(kw)}

		group, err := l.lexGroupToken(false)
		if err != nil {
			return nil, fmt.Errorf("error in keyword '%s' at position %d: %w", kw, l.Pos(), err)
		}

		tokens = append(tokens, group)
	}

	return tokens, nil
}

func MustNewRexTokensFromKeywords(kws ...string) []dialect.Token {
	tokens, err := NewRexTokensFromKeywords(kws...)
	if err != nil {
		panic(err)
	}

	return tokens
}

type keywordsLexer struct {
	lexer.Lexer
}

var (
	ErrEOF            = errors.New("EOF")
	ErrUnexpectedEOF  = errors.New("unexpected EOF")
	ErrUnexpectedChar = errors.New("unexpected character")
)

func (l *keywordsLexer) lexGroupToken(parens bool) (base.GroupToken, error) {
	var groupTokens []dialect.Token
outer:
	for {
		var tokens []dialect.Token
		addGroup := func() {
			if len(tokens) > 0 {
				groupTokens = append(groupTokens, rex.Group.NonCaptured(tokens...))
			}
			tokens = nil
		}
	inner:
		for {
			if parens {
				if l.ReadChar('(') {
					l.Backup()
					return base.GroupToken{}, ErrUnexpectedChar
				}
				if l.ReadChar(')') {
					addGroup()
					break outer
				}
			} else if l.ReadChar('(') {
				group, err := l.lexGroupToken(true)
				if err != nil {
					return base.GroupToken{}, err
				}
				if l.ReadChar('?') {
					tokens = append(tokens, group.Repeat().ZeroOrOne())
				} else {
					tokens = append(tokens, group)
				}
				continue inner
			}
			if l.ReadChar('|') {
				if len(tokens) == 0 {
					l.Backup()
					return base.GroupToken{}, ErrUnexpectedChar
				}
				addGroup()
				continue outer
			}
			token, err := l.lexClassWithModifierToken()
			if errors.Is(err, ErrEOF) {
				if parens {
					return base.GroupToken{}, ErrUnexpectedEOF
				}
				addGroup()
				break outer
			}
			if err != nil {
				return base.GroupToken{}, err
			}
			tokens = append(tokens, token)
			continue inner
		}
	}

	if len(groupTokens) == 0 {
		return base.GroupToken{}, ErrUnexpectedEOF
	}

	return rex.Group.Composite(groupTokens...).NonCaptured(), nil
}

func (l *keywordsLexer) lexClassWithModifierToken() (dialect.Token, error) {
	if l.ReadChar('*') {
		return regex.AnyWordChar().Repeat().ZeroOrMore(), nil
	}

	tk, err := l.lexClassToken()
	if err == nil {
		ch, ok := l.Read()

		switch {
		case !ok:
			return tk, nil
		case ch == '?':
			return tk.Repeat().ZeroOrOne(), nil
		case ch == '+':
			return tk.Repeat().OneOrMore(), nil
		default:
			l.Backup()
			return tk, nil
		}
	}

	return nil, err
}

var reservedChars = map[rune]struct{}{
	'(': {}, ')': {}, '|': {}, '*': {}, '?': {}, '+': {}, '#': {}, ' ': {},
}

func (l *keywordsLexer) lexClassToken() (base.ClassToken, error) {
	var tk base.ClassToken

	ch, ok := l.Read()

	switch {
	case !ok:
		return tk, ErrEOF
	case ch == '\\':
		exactChar, ok2 := l.Read()
		if !ok2 {
			return tk, ErrUnexpectedEOF
		}

		return rex.Chars.Single(exactChar), nil
	case lexer.IsWordChar(ch):
		lcChar := strings.ToLower(string(ch))
		if string(ch) != lcChar {
			return tk, ErrUnexpectedChar
		}

		ucChar := strings.ToUpper(string(ch))
		if lcChar == ucChar {
			tk = rex.Chars.Single(ch)
		} else {
			tk = rex.Chars.Runes(ucChar + lcChar)
		}

		return tk, nil
	case ch == '#':
		return rex.Chars.Digits(), nil
	case ch == ' ':
		return regex.AnyNonWordChar(), nil
	default:
		if _, ok := reservedChars[ch]; ok {
			l.Backup()
			return tk, ErrUnexpectedChar
		}

		return rex.Chars.Single(ch), nil
	}
}
