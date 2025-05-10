package lexer

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func NewLexer(str string) Lexer {
	return Lexer{
		reader: bufio.NewReader(strings.NewReader(str)),
	}
}

type Lexer struct {
	pos    int
	reader *bufio.Reader
}

func (l *Lexer) Pos() int {
	return l.pos
}

func (l *Lexer) Read() (rune, bool) {
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

func (l *Lexer) Backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos--
}

func (l *Lexer) BackupN(n int) {
	for range n {
		l.Backup()
	}
}

func (l *Lexer) IsEOF() bool {
	_, ok := l.Read()
	if !ok {
		return true
	}

	l.Backup()

	return false
}

func (l *Lexer) ReadIf(fn func(rune) bool) (rune, bool) {
	r, ok := l.Read()
	if !ok {
		return 0, false
	}

	if !fn(r) {
		l.Backup()
		return 0, false
	}

	return r, true
}

func (l *Lexer) ReadWhile(fn func(rune) bool) string {
	var str string

	for {
		r, ok := l.ReadIf(fn)
		if !ok {
			break
		}

		str += string(r)
	}

	return str
}

func (l *Lexer) ReadUntil(fn func(rune) bool) string {
	var str string

	for {
		r, ok := l.Read()
		if !ok {
			break
		}

		str += string(r)

		if fn(r) {
			break
		}
	}

	return str
}

func (l *Lexer) ReadInt() (int, bool) {
	str := l.ReadWhile(IsInt)
	if str == "" {
		return 0, false
	}

	n, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return n, true
}

func (l *Lexer) ReadChar(r1 rune) bool {
	_, ok := l.ReadIf(IsChar(r1))
	return ok
}

func IsInt(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsChar(r1 rune) func(rune) bool {
	return func(r2 rune) bool {
		return r1 == r2
	}
}

func IsWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func IsNonWordChar(r rune) bool {
	return !IsWordChar(r)
}
