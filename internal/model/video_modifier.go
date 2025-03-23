package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/keywords"
)

// VideoModifier represents the modifier of a video
// ENUM(REGIONAL, SCREENER, RAWHD, BRDISK, REMUX)
type VideoModifier string

func (v VideoModifier) Label() string {
	return v.String()
}

var videoModifierRegex = keywords.MustNewRegexFromKeywords(namesToLower(VideoModifierNames()...)...)

func InferVideoModifier(input string) NullVideoModifier {
	if match := videoModifierRegex.FindStringSubmatch(input); match != nil {
		if parsed, parseErr := ParseVideoModifier(match[1]); parseErr == nil {
			return NewNullVideoModifier(parsed)
		}
	}

	return NullVideoModifier{}
}
