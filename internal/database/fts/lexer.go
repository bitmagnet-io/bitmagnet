package fts

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func newLexer(str string) lexer {
	return lexer{
		reader: bufio.NewReader(strings.NewReader(str)),
	}
}

type lexer struct {
	pos    int
	reader *bufio.Reader
}

func (l *lexer) isEof() bool {
	_, ok := l.read()
	if !ok {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) read() (rune, bool) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, false
		}
		panic(err)
	}
	l.pos++
	return r, true
}

func (l *lexer) readIf(fn func(rune) bool) (rune, bool) {
	r, ok := l.read()
	if !ok {
		return 0, false
	}
	if !fn(r) {
		l.backup()
		return 0, false
	}
	return r, true
}

func (l *lexer) readWhile(fn func(rune) bool) string {
	var str string
	for {
		r, ok := l.readIf(fn)
		if !ok {
			break
		}
		str = str + string(r)
	}
	return str
}

func (l *lexer) readInt() (int, bool) {
	str := l.readWhile(isInt)
	if str == "" {
		return 0, false
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return n, true
}

func (l *lexer) readChar(r1 rune) bool {
	_, ok := l.readIf(isChar(r1))
	return ok
}

func (l *lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.pos--
}

func (l *lexer) readQuotedString(quoteChar rune) (string, error) {
	if !l.readChar(quoteChar) {
		return "", errors.New("missing opening quote")
	}
	var str string
	for {
		ch, ok := l.read()
		if !ok {
			return str, errors.New("unexpected EOF")
		}
		if ch == quoteChar && !l.readChar(quoteChar) {
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
