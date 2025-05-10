package fts

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
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

		str += string(ch)
	}

	return str, nil
}
