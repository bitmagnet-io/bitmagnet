package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"regexp"
	"strings"
)

// VideoResolution represents the resolution of a video
// ENUM(V360p, V480p, V540p, V576p, V720p, V1080p, V1440p, V2160p, V4320p)
type VideoResolution string

func (v VideoResolution) Label() string {
	return v.String()[1:]
}

var videoResolutionAliases = map[string]VideoResolution{
	"1080i":     VideoResolutionV1080p,
	"1920x1080": VideoResolutionV1080p,
	"3840x2160": VideoResolutionV2160p,
	"8k":        VideoResolutionV4320p,
	"4k":        VideoResolutionV2160p,
	"uhd":       VideoResolutionV2160p,
	"2k":        VideoResolutionV1440p,
	"hd":        VideoResolutionV1080p,
	"sd":        VideoResolutionV480p,
}

func createVideoResolutionRegex() *regexp.Regexp {
	names := removeEnumPrefixes(VideoResolutionNames()...)
	for res := range videoResolutionAliases {
		names = append(names, res)
	}
	return regex.NewRegexFromNames(names...)
}

var videoResolutionRegex = createVideoResolutionRegex()

func InferVideoResolution(input string) NullVideoResolution {
	if match := videoResolutionRegex.FindStringSubmatch(input); match != nil {
		if parsed, parseErr := ParseVideoResolution("V" + match[1]); parseErr == nil {
			return NewNullVideoResolution(parsed)
		}
		lowerMatch := strings.ToLower(match[1])
		if inferred, ok := videoResolutionAliases[lowerMatch]; ok {
			return NewNullVideoResolution(inferred)
		}
	}
	return NullVideoResolution{}
}
