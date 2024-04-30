package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/keywords"
)

// Video3d represents the 3D type of a video
// ENUM(V3D, V3DSBS, V3DOU)
type Video3d string

func (v Video3d) Label() string {
	return v.String()[1:]
}

var video3dRegex = keywords.MustNewRegexFromKeywords(namesToLower(removeEnumPrefixes(Video3dNames()...)...)...)

func InferVideo3d(input string) NullVideo3d {
	if match := video3dRegex.FindStringSubmatch(input); match != nil {
		if parsed, parseErr := ParseVideo3d("V" + match[1]); parseErr == nil {
			return NewNullVideo3d(parsed)
		}
	}
	return NullVideo3d{}
}
