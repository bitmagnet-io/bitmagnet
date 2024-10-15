package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/keywords"
)

// Video3D represents the 3D type of a video
// ENUM(V3D, V3DSBS, V3DOU)
type Video3D string

func (v Video3D) Label() string {
	return v.String()[1:]
}

var video3dRegex = keywords.MustNewRegexFromKeywords(namesToLower(removeEnumPrefixes(Video3DNames()...)...)...)

func InferVideo3D(input string) NullVideo3D {
	if match := video3dRegex.FindStringSubmatch(input); match != nil {
		if parsed, parseErr := ParseVideo3D("V" + match[1]); parseErr == nil {
			return NewNullVideo3D(parsed)
		}
	}
	return NullVideo3D{}
}
