package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/model"

type ContentAttributes struct {
	Languages       model.Languages
	LanguageMulti   bool
	Episodes        model.Episodes
	VideoResolution model.NullVideoResolution
	VideoSource     model.NullVideoSource
	VideoCodec      model.NullVideoCodec
	Video3d         model.NullVideo3d
	VideoModifier   model.NullVideoModifier
	ReleaseGroup    model.NullString
}

type Classification struct {
	ContentType model.NullContentType
	Content     *model.Content
	ContentAttributes
}

func (a *ContentAttributes) ApplyHint(h model.TorrentHint) {
	if len(h.Episodes) > 0 {
		a.Episodes = h.Episodes
	}
	if len(h.Languages) > 0 {
		a.Languages = h.Languages
	}
	if h.VideoResolution.Valid {
		a.VideoResolution = h.VideoResolution
	}
	if h.VideoSource.Valid {
		a.VideoSource = h.VideoSource
	}
	if h.VideoCodec.Valid {
		a.VideoCodec = h.VideoCodec
	}
	if h.Video3d.Valid {
		a.Video3d = h.Video3d
	}
	if h.VideoModifier.Valid {
		a.VideoModifier = h.VideoModifier
	}
	if h.ReleaseGroup.Valid {
		a.ReleaseGroup = h.ReleaseGroup
	}
}

func (c *Classification) ApplyHint(h model.TorrentHint) {
	c.ContentType = h.NullContentType()
	c.ContentAttributes.ApplyHint(h)
}
