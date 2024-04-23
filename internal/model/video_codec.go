package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/keywords"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/hedhyw/rex/pkg/rex"
	"regexp"
	"strings"
)

// VideoCodec represents the codec of a video
// ENUM(H264, x264, x265, XviD, DivX, MPEG2, MPEG4)
type VideoCodec string

func (v VideoCodec) Label() string {
	return v.String()
}

var videoCodecAliases = map[string]VideoCodec{
	"avc": VideoCodecH264,
}

func createVideoCodecAndOptionalReleaseGroupRegex() *regexp.Regexp {
	names := namesToLower(VideoCodecNames()...)
	for videoCodec := range videoCodecAliases {
		names = append(names, videoCodec)
	}
	return rex.New(
		rex.Group.Composite(
			rex.Chars.Begin(),
			regex.AnyNonWordChar().Repeat().OneOrMore(),
		).NonCaptured(),
		rex.Group.Composite(
			keywords.MustNewRexTokensFromKeywords(names...)...,
		),
		rex.Group.Composite(
			rex.Chars.End(),
			rex.Group.NonCaptured(
				rex.Chars.Single('-'),
				rex.Group.Define(regex.AnyWordChar().Repeat().OneOrMore()),
			),
			rex.Common.NotClass(regex.AnyNonWordChar()).Repeat().OneOrMore(),
		).NonCaptured(),
	).MustCompile()
}

var videoCodecAndOptionalReleaseGroupRegex = createVideoCodecAndOptionalReleaseGroupRegex()

func InferVideoCodecAndReleaseGroup(input string) (NullVideoCodec, NullString) {
	if match := videoCodecAndOptionalReleaseGroupRegex.FindStringSubmatch(input); match != nil {
		releaseGroup := NullString{
			String: match[2],
			Valid:  match[2] != "",
		}
		if videoCodec, parseErr := ParseVideoCodec(match[1]); parseErr == nil {
			return NewNullVideoCodec(videoCodec), releaseGroup
		}
		if aliased, ok := videoCodecAliases[strings.ToLower(match[1])]; ok {
			return NewNullVideoCodec(aliased), releaseGroup
		}
	}
	return NullVideoCodec{}, NullString{}
}
