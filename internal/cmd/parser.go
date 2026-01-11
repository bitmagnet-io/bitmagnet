package cmd

import (
	"regexp"
	"strings"
)

type tokenType int

const (
	tokenTypeUnknown tokenType = iota
	tokenTypeKeyValue
	tokenTypeShortKeyValue
	tokenTypeKey
	tokenTypeShortKey
	tokenTypePositional
)

type token struct {
	tokenType
	arg   string
	key   string
	value string
}

func parseToken(arg string) token {
	if key, value, ok := parseKeyValue("--", arg); ok {
		return token{
			tokenType: tokenTypeKeyValue,
			arg:       arg,
			key:       key,
			value:     value,
		}
	}

	if key, ok := parseKey("--", arg); ok {
		return token{
			tokenType: tokenTypeKey,
			arg:       arg,
			key:       key,
		}
	}

	if key, value, ok := parseKeyValue("-", arg); ok {
		return token{
			tokenType: tokenTypeShortKeyValue,
			arg:       arg,
			key:       key,
			value:     value,
		}
	}

	if key, ok := parseKey("-", arg); ok {
		return token{
			tokenType: tokenTypeShortKey,
			arg:       arg,
			key:       key,
		}
	}

	return token{
		tokenType: tokenTypePositional,
		arg:       arg,
		value:     arg,
	}
}

var (
	regexParamKey      = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	regexParamKeyValue = regexp.MustCompile(`^([a-z0-9]+(?:-[a-z0-9]+)*)=(.+)$`)
)

func parseKeyValue(prefix, arg string) (string, string, bool) {
	if !strings.HasPrefix(arg, prefix) || len(arg) <= len(prefix) {
		return "", "", false
	}

	match := regexParamKeyValue.FindStringSubmatch(arg[2:])

	if len(match) == 0 {
		return "", "", false
	}

	return match[1], match[2], true
}

func parseKey(prefix, arg string) (string, bool) {
	if !strings.HasPrefix(arg, prefix) || len(arg) <= len(prefix) {
		return "", false
	}

	match := regexParamKey.FindStringSubmatch(arg[2:])

	if len(match) == 0 {
		return "", false
	}

	return match[0], true
}
