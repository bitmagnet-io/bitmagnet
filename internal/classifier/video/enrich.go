package video

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func PreEnrich(input model.TorrentContent) (model.TorrentContent, error) {
	if hasVideo := input.Torrent.HasFileType(model.FileTypeVideo); hasVideo.Valid && !hasVideo.Bool {
		return model.TorrentContent{}, classifier.ErrNoMatch
	}
	if input.ContentType.Valid && !input.ContentType.ContentType.IsVideo() {
		return model.TorrentContent{}, classifier.ErrNoMatch
	}
	core, rest, coreInfoErr := ParseVideoCoreInfo(input.ContentType, input.Torrent.Name)
	if coreInfoErr != nil {
		return model.TorrentContent{}, coreInfoErr
	}
	output := input
	if len(output.Episodes) == 0 {
		output.Episodes = core.Episodes
	}
	if !output.ContentType.Valid {
		if len(output.Episodes) > 0 {
			output.ContentType.ContentType = model.ContentTypeTvShow
		} else {
			output.ContentType.ContentType = model.ContentTypeMovie
		}
		output.ContentType.Valid = true
	}
	output.Title = core.Title
	if output.ReleaseYear.IsNil() {
		output.ReleaseYear = core.Year
	}
	if len(output.Languages) == 0 {
		langs := model.InferLanguages(rest)
		if len(langs) > 0 {
			output.Languages = langs
		}
	}
	if !output.VideoResolution.Valid {
		output.VideoResolution = model.InferVideoResolution(rest)
	}
	if !output.VideoSource.Valid {
		output.VideoSource = model.InferVideoSource(rest)
	}
	if !output.VideoModifier.Valid {
		output.VideoModifier = model.InferVideoModifier(rest)
	}
	if !output.Video3d.Valid {
		output.Video3d = model.InferVideo3d(rest)
	}
	if !output.VideoCodec.Valid || !output.ReleaseGroup.Valid {
		vc, rg := model.InferVideoCodecAndReleaseGroup(rest)
		if !output.VideoCodec.Valid {
			output.VideoCodec = vc
		}
		if !output.ReleaseGroup.Valid {
			output.ReleaseGroup = rg
		}
	}
	return output, nil
}
