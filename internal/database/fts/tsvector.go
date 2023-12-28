package fts

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type TsvectorWeight rune

const (
	TsvectorWeightA TsvectorWeight = 'A'
	TsvectorWeightB TsvectorWeight = 'B'
	TsvectorWeightC TsvectorWeight = 'C'
	TsvectorWeightD TsvectorWeight = 'D'
)

type TsvectorLabel struct {
	Position int
	Weight   TsvectorWeight
}

type Tsvector map[string]map[TsvectorLabel]struct{}

func (v Tsvector) String() string {
	var entries []maps.MapEntry[string, []TsvectorLabel]
	for lexeme, labelsMap := range v {
		labels := make([]TsvectorLabel, 0, len(labelsMap))
		for label := range labelsMap {
			labels = append(labels, label)
		}
		sort.Slice(labels, func(i, j int) bool {
			return labels[i].Position < labels[j].Position
		})
		entries = append(entries, maps.MapEntry[string, []TsvectorLabel]{
			Key:   lexeme,
			Value: labels,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return ""
}

func ParseTsvector(str string) (Tsvector, error) {
	l := tsvectorLexer{newLexer(strings.TrimSpace(str))}
	tsv := NewTsvector()
	for {
		if l.isEof() {
			break
		}
		word, posWeights, err := l.readTsvPart()
		if err != nil {
			return nil, fmt.Errorf("error at position %d: %w", l.pos, err)
		}
		if _, ok := tsv[word]; !ok {
			tsv[word] = make(map[TsvectorLabel]struct{})
		}
		for _, posWeight := range posWeights {
			tsv[word][posWeight] = struct{}{}
		}
	}
	return tsv, nil
}

type tsvectorLexer struct {
	lexer
}

func (l *tsvectorLexer) readTsvPart() (string, []TsvectorLabel, error) {
	lexeme, err := l.readLexeme()
	if err != nil {
		return "", nil, err
	}
	if !l.readChar(':') {
		spaces := l.readWhile(isChar(' '))
		if !l.isEof() && len(spaces) == 0 {
			return "", nil, errors.New("unexpected character")
		}
		return lexeme, []TsvectorLabel{{1, TsvectorWeightD}}, nil
	}
	labels, err := l.readLabels()
	if err != nil {
		return "", nil, err
	}
	return lexeme, labels, nil
}

func (l *tsvectorLexer) readLexeme() (string, error) {
	word, err := l.readQuotedString('\'')
	if err != nil {
		return "", err
	}
	if word == "" {
		return "", errors.New("empty quoted string")
	}
	return word, nil
}

func (l *tsvectorLexer) readLabels() ([]TsvectorLabel, error) {
	pos, ok := l.readInt()
	if !ok {
		return nil, errors.New("missing position")
	}
	pw := TsvectorLabel{
		Position: pos,
	}
	if w, ok := l.readIf(isWeight); ok {
		pw.Weight = TsvectorWeight(unicode.ToUpper(w))
	}
	pws := []TsvectorLabel{pw}
	if l.readChar(',') {
		if rest, err := l.readLabels(); err != nil {
			return nil, err
		} else {
			pws = append(pws, rest...)
		}
	} else if !l.isEof() && !l.readChar(' ') {
		return nil, errors.New("unexpected character")
	}
	return pws, nil
}

func isWeight(r rune) bool {
	r = unicode.ToUpper(r)
	return r == 'A' || r == 'B' || r == 'C' || r == 'D'
}

func NewTsvector() Tsvector {
	return Tsvector{}
}

func (t Tsvector) Add(label TsvectorWeight, value string) {
	nextPos := 1
	for _, pls := range t {
		for pl := range pls {
			if pl.Weight == label && pl.Position > nextPos {
				nextPos = pl.Position + 1
			}
		}
	}
	for _, lexeme := range TokenizeFlat(value) {
		if _, ok := t[lexeme]; !ok {
			t[lexeme] = make(map[TsvectorLabel]struct{})
		}
		t[lexeme][TsvectorLabel{
			Position: nextPos,
			Weight:   label,
		}] = struct{}{}
		nextPos++
	}
}

var nonWordChar = regexp.MustCompile(`\W+`)

func quoteLexeme(str string) string {
	if nonWordChar.MatchString(str) {
		str = strings.Replace(str, "\\", "\\\\", -1)
		str = "'" + strings.Replace(str, "'", "\\'", -1) + "'"
	}
	return str
}
