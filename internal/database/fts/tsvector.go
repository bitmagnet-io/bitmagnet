package fts

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

type Tsvector map[string]map[int]TsvectorWeight

func (v Tsvector) Copy() Tsvector {
	c := Tsvector{}
	for lexeme, labels := range v {
		c[lexeme] = make(map[int]TsvectorWeight)
		for pos, weight := range labels {
			c[lexeme][pos] = weight
		}
	}

	return c
}

func (v Tsvector) String() string {
	entries := make([]maps.MapEntry[string, []TsvectorLabel], 0, len(v))

	for lexeme, labelsMap := range v {
		n := len(labelsMap)
		if n == 0 {
			n = 1
		}

		labels := make([]TsvectorLabel, 0, n)
		if len(labelsMap) == 0 {
			labels = append(labels, TsvectorLabel{0, TsvectorWeightD})
		} else {
			for pos, weight := range labelsMap {
				labels = append(labels, TsvectorLabel{pos, weight})
			}

			sort.Slice(labels, func(i, j int) bool {
				return labels[i].Position < labels[j].Position
			})
		}

		entries = append(entries, maps.MapEntry[string, []TsvectorLabel]{
			Key:   lexeme,
			Value: labels,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	parts := make([]string, 0, len(entries))

	for _, entry := range entries {
		var labels []string

		for _, l := range entry.Value {
			if l.Position == 0 {
				continue
			}

			label := fmt.Sprintf("%d", l.Position)
			if l.Weight != TsvectorWeightD {
				label += string(l.Weight)
			}

			labels = append(labels, label)
		}

		part := quoteLexeme(entry.Key, true)
		if len(labels) > 0 {
			part = fmt.Sprintf("%s:%s", part, strings.Join(labels, ","))
		}

		parts = append(parts, part)
	}

	return strings.Join(parts, " ")
}

func ParseTsvector(str string) (Tsvector, error) {
	l := tsvectorLexer{newLexer(strings.TrimSpace(str))}
	tsv := Tsvector{}

	for !l.IsEOF() {
		word, posWeights, err := l.readTsvPart()
		if err != nil {
			return nil, fmt.Errorf("error at position %d: %w", l.Pos(), err)
		}

		if _, ok := tsv[word]; !ok {
			tsv[word] = make(map[int]TsvectorWeight)
		}

		for _, posWeight := range posWeights {
			tsv[word][posWeight.Position] = posWeight.Weight
		}
	}

	return tsv, nil
}

type tsvectorLexer struct {
	ftsLexer
}

func (l *tsvectorLexer) readTsvPart() (string, []TsvectorLabel, error) {
	lexeme, err := l.readLexeme()
	if err != nil {
		return "", nil, err
	}

	if !l.ReadChar(':') {
		spaces := l.ReadWhile(lexer.IsChar(' '))
		if !l.IsEOF() && len(spaces) == 0 {
			return "", nil, errors.New("unexpected character")
		}

		return lexeme, []TsvectorLabel{}, nil
	}

	labels, err := l.readLabels()
	if err != nil {
		return "", nil, err
	}

	return lexeme, labels, nil
}

func (l *tsvectorLexer) readLexeme() (string, error) {
	if unquoted := l.ReadWhile(lexer.IsWordChar); unquoted != "" {
		return unquoted, nil
	}

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
	pos, ok := l.ReadInt()
	if !ok {
		return nil, errors.New("missing position")
	}

	pw := TsvectorLabel{
		Position: pos,
		Weight:   TsvectorWeightD,
	}
	if w, ok := l.ReadIf(isWeight); ok {
		pw.Weight = TsvectorWeight(unicode.ToUpper(w))
	}

	pws := []TsvectorLabel{pw}

	if l.ReadChar(',') {
		rest, err := l.readLabels()
		if err != nil {
			return nil, err
		}

		pws = append(pws, rest...)
	} else if !l.IsEOF() && l.ReadWhile(lexer.IsChar(' ')) == "" {
		return nil, errors.New("unexpected character")
	}

	return pws, nil
}

func isWeight(r rune) bool {
	r = unicode.ToUpper(r)
	return r == 'A' || r == 'B' || r == 'C' || r == 'D'
}

func (v Tsvector) AddText(text string, weight TsvectorWeight) {
	nextPos := 1

	for _, pls := range v {
		for pos := range pls {
			if pos >= nextPos {
				nextPos = pos + 1
			}
		}
	}

	if nextPos > 1 {
		nextPos++
	}

	for _, lexeme := range TokenizeFlat(text) {
		if _, ok := v[lexeme]; !ok {
			v[lexeme] = make(map[int]TsvectorWeight)
		}

		v[lexeme][nextPos] = weight
		nextPos++
	}
}

var nonWordChar = regexp.MustCompile(`\W`)

func quoteLexeme(str string, force bool) string {
	if force || nonWordChar.MatchString(str) {
		str = "'" + strings.ReplaceAll(str, "'", "''") + "'"
	}

	return str
}

func (v *Tsvector) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	str, ok := val.(string)
	if !ok {
		return errors.New("invalid type")
	}

	parsed, err := ParseTsvector(str)
	if err != nil {
		return err
	}

	*v = parsed

	return nil
}

func (Tsvector) Value() (driver.Value, error) {
	return nil, errors.New("cannot get value")
}

func (Tsvector) GormDataType() string {
	return "tsvector"
}

func (v Tsvector) GormValue(context.Context, *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "?::tsvector",
		Vars: []interface{}{v.String()},
	}
}
