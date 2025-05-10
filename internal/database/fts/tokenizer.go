package fts

import (
	"strings"
	"unicode"

	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
	"github.com/mozillazg/go-unidecode/table"
)

func Tokenize(str string) [][]string {
	l := tokenizerLexer{newLexer(str)}

	var tokens [][]string

	for {
		phrase := l.readPhrase()
		if len(phrase) == 0 {
			break
		}

		tokens = append(tokens, phrase)
	}

	return tokens
}

type tokenizerLexer struct {
	ftsLexer
}

func TokenizeFlat(str string) []string {
	var tokens []string
	for _, phrase := range Tokenize(str) {
		tokens = append(tokens, phrase...)
	}

	return tokens
}

func (l *tokenizerLexer) readPhrase() []string {
	var phrase []string

	var lexeme string

	breakWord := func() {
		if lexeme != "" {
			phrase = append(phrase, lexeme)
			lexeme = ""
		}
	}
	appendStr := func(str string) {
		lexeme += str
	}

	for {
		if l.IsEOF() {
			breakWord()
			return phrase
		}

		if ch, ok := l.ReadIf(lexer.IsWordChar); ok {
			ch = unicode.ToLower(ch)
			if ch < unicode.MaxASCII {
				appendStr(string(ch))
			} else {
				// If the character is determined to be a language with unspaced words
				// (e.g. Chinese, Japanese), each character will become a token;
				// using this cutoff might not be perfect.
				isNonBreakingLang := ch > '\u1FFF'
				if isNonBreakingLang {
					breakWord()
				}

				section := ch >> 8   // Chop off the last two hex digits
				position := ch % 256 // Last two hex digits

				if tb, ok := table.Tables[section]; ok {
					if len(tb) > int(position) {
						subst := tb[position]
						// replace some problematic characters
						subst = strings.ReplaceAll(subst, "'", "_sq_")
						subst = strings.ReplaceAll(subst, "\\", "_bs_")
						subst = strings.TrimSpace(subst)
						appendStr(subst)

						if isNonBreakingLang || len(subst) == 0 || subst[len(subst)-1] == ' ' {
							breakWord()
						}
					}
				}
			}

			continue
		}

		breakWord()

		if len(phrase) > 0 {
			return phrase
		}

		l.Read()
	}
}
