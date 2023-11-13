package adult

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func PreEnrich(input model.TorrentContent) (model.TorrentContent, error) {
	if !strings.Contains(strings.ToLower(input.Title), "xxx") {
		return model.TorrentContent{}, resolver.ErrNoMatch
	}

	titleLower := strings.ToLower(input.Title)
	titleLower = strings.Replace(titleLower, "com_", "", 0)
	titleLower = strings.Replace(titleLower, "www.torrenting.com", "", 0)
	titleLower = strings.Replace(titleLower, "www.torrenting.org", "", 0)

	output := input
	output.Title = titleLower

	if !output.VideoResolution.Valid {
		output.VideoResolution = model.InferVideoResolution(output.Title)
	}
	if !output.VideoSource.Valid {
		output.VideoSource = model.InferVideoSource(output.Title)
	}
	if !output.VideoModifier.Valid {
		output.VideoModifier = model.InferVideoModifier(output.Title)
	}
	if !output.Video3d.Valid {
		output.Video3d = model.InferVideo3d(output.Title)
	}
	if !output.VideoCodec.Valid || !output.ReleaseGroup.Valid {
		vc, rg := model.InferVideoCodecAndReleaseGroup(output.Title)
		if !output.VideoCodec.Valid {
			output.VideoCodec = vc
		}
		if !output.ReleaseGroup.Valid {
			output.ReleaseGroup = rg
		}
	}
	return output, nil
}
