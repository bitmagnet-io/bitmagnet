package cli

import (
	"bytes"
	_ "embed"
	"io"
	"regexp"

	"github.com/bitmagnet-io/bitmagnet/internal/colors"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/fatih/color"
)

var (
	//go:embed banner.txt
	rawBanner []byte

	patternBlock = regexp.MustCompile("([" + blockChars + "]+)")
	patternText  = regexp.MustCompile("([^" + blockChars + "v ]+)")
)

const blockChars = "█▀▄▝▜▛▘▙▟▂▁"

type transform func(bytes []byte) []byte

func writeBanner(writer io.Writer) {
	_, _ = writer.Write([]byte("\n"))
	_, _ = writer.Write(
		transforms(
			replaceColor(patternText, colors.White),
			replaceColor(patternBlock, colors.Blue),
			replaceVersion,
		)(rawBanner),
	)
	_, _ = writer.Write([]byte("\n"))
}

func transforms(transforms ...transform) transform {
	return func(bs []byte) []byte {
		for _, t := range transforms {
			bs = t(bs)
		}

		return bs
	}
}

func replaceColor(pattern *regexp.Regexp, color *color.Color) transform {
	return func(bs []byte) []byte {
		return pattern.ReplaceAllFunc(
			bs,
			func(m []byte) []byte {
				return []byte(color.Sprintf("%s", string(m)))
			},
		)
	}
}

func replaceVersion(bs []byte) []byte {
	return bytes.Replace(bs, []byte("v"), []byte(colors.Blue.Sprintf("%s", version.GitTag)), 1)
}
