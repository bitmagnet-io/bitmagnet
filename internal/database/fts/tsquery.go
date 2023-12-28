package fts

import (
	"strings"
)

func AppQueryToTsquery(str string) string {
	l := queryLexer{newLexer(str)}
	var tokens []TokenValue
	for {
		token, ok := l.readQueryToken()
		if !ok {
			break
		}
		tokens = append(tokens, token)
	}
	return appQueryTokensToTsquery(tokens...)
}

type queryLexer struct {
	lexer
}

type Token string

const (
	TokenQuoted      Token = "QUOTED"
	TokenOpenParens  Token = "OPEN_PARENS"
	TokenCloseParens Token = "CLOSE_PARENS"
	TokenOperator    Token = "OPERATOR"
	TokenNegation    Token = "NEGATION"
	TokenPhrase      Token = "PHRASE"
)

type TokenValue struct {
	StartPos int
	Length   int
	Token    Token
	Value    string
}

type Operator string

const (
	OperatorAnd        Operator = "&"
	OperatorOr         Operator = "|"
	OperatorFollowedBy Operator = "."
)

func (l *queryLexer) readQueryToken() (TokenValue, bool) {
	for {
		if l.isEof() {
			return TokenValue{}, false
		}
		start := l.pos
		if l.readChar('(') {
			return TokenValue{l.pos, 1, TokenOpenParens, "("}, true
		}
		if l.readChar(')') {
			return TokenValue{l.pos, 1, TokenCloseParens, ")"}, true
		}
		if l.readChar('&') {
			return TokenValue{l.pos, 1, TokenOperator, string(OperatorAnd)}, true
		}
		if l.readChar('|') {
			return TokenValue{l.pos, 1, TokenOperator, string(OperatorOr)}, true
		}
		if l.readChar('.') {
			return TokenValue{l.pos, 1, TokenOperator, string(OperatorFollowedBy)}, true
		}
		if l.readChar('!') {
			return TokenValue{l.pos, 1, TokenNegation, "!"}, true
		}
		if quoted, _ := l.readQuotedString('"'); quoted != "" {
			return TokenValue{start, l.pos - start, TokenQuoted, quoted}, true
		}
		if phrase := l.readWhile(isWordChar); phrase != "" {
			return TokenValue{start, l.pos - start, TokenPhrase, phrase}, true
		}
		l.read()
	}
}

func appQueryTokensToTsquery(tokens ...TokenValue) string {
	var parts []string
	i := 0
outer:
	for {
		var operator Operator
		var negated bool
		for {
			if i >= len(tokens) {
				break outer
			}
			token := tokens[i]
			addExpr := func(expr string) {
				if len(parts) > 0 && operator == "" {
					parts = append(parts, string(OperatorAnd))
				}
				if operator != "" {
					strOp := string(operator)
					if operator == OperatorFollowedBy {
						strOp = "<->"
					}
					parts = append(parts, strOp)
				}
				if negated {
					parts = append(parts, "!")
				}
				parts = append(parts, expr)
				operator = ""
				negated = false
			}
			switch token.Token {
			case TokenOperator:
				operator = Operator(token.Value)
			case TokenNegation:
				negated = true
			case TokenQuoted:
				tokenized := TokenizeFlat(token.Value)
				var quotedWords []string
				for _, word := range tokenized {
					quotedWords = append(quotedWords, quoteLexeme(word))
				}
				if len(quotedWords) > 0 {
					addExpr(strings.Join(quotedWords, " <-> "))
				}
			case TokenPhrase:
				tokenized := Tokenize(token.Value)
				var phrases []string
				for _, phrase := range tokenized {
					var quotedWords []string
					for _, word := range phrase {
						quotedWords = append(quotedWords, quoteLexeme(word))
					}
					if len(quotedWords) > 0 {
						phrases = append(phrases, strings.Join(quotedWords, " <-> "))
					}
				}
				if len(phrases) > 0 {
					addExpr(strings.Join(phrases, " & "))
				}
			case TokenOpenParens:
				var parensTokens []TokenValue
				depth := 1
				for j := i + 1; j < len(tokens); j++ {
					if tokens[j].Token == TokenOpenParens {
						depth++
					}
					if tokens[j].Token == TokenCloseParens {
						depth--
						if depth == 0 {
							break
						}
					}
					parensTokens = append(parensTokens, tokens[j])
				}
				i += len(parensTokens)
				parensExpr := appQueryTokensToTsquery(parensTokens...)
				if len(parensExpr) > 0 {
					addExpr("(" + parensExpr + ")")
				}
			}
			i++
		}
	}
	return strings.Join(parts, " ")
}
