package fts

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
	"unicode"
)

func newLexer(str string) ftsLexer {
	return ftsLexer{lexer.NewLexer(str)}
}

type ftsLexer struct {
	lexer.Lexer
}

func (l *ftsLexer) readQuotedString(quoteChar rune) (string, error) {
	if !l.ReadChar(quoteChar) {
		return "", errors.New("missing opening quote")
	}
	var str string
	for {
		ch, ok := l.Read()
		if !ok {
			return str, errors.New("unexpected EOF")
		}
		if ch == quoteChar && !l.ReadChar(quoteChar) {
			break
		}
		str = str + string(ch)
	}
	return str, nil
}

func isInt(r rune) bool {
	return r >= '0' && r <= '9'
}

func isChar(r1 rune) func(rune) bool {
	return func(r2 rune) bool {
		return r1 == r2
	}
}

func IsWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
