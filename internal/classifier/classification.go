package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/model"

type Classification struct {
	ContentAttributes
	Content *model.Content
}

func (c *Classification) ApplyHint(h model.TorrentHint) {
	c.ContentType = h.NullContentType()
	c.ContentAttributes.ApplyHint(h)
}

type ContentAttributes struct {
	ContentType     model.NullContentType
	BaseTitle       model.NullString
	Year            model.Year
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

func (a *ContentAttributes) Merge(other ContentAttributes) {
	if !a.ContentType.Valid {
		a.ContentType = other.ContentType
	}
	if !a.BaseTitle.Valid {
		a.BaseTitle = other.BaseTitle
	}
	if a.Year.IsNil() {
		a.Year = other.Year
	}
	if len(a.Languages) == 0 {
		a.Languages = other.Languages
	}
	a.LanguageMulti = a.LanguageMulti || other.LanguageMulti
	if len(a.Episodes) == 0 {
		a.Episodes = other.Episodes
	}
	if !a.VideoResolution.Valid {
		a.VideoResolution = other.VideoResolution
	}
	if !a.VideoSource.Valid {
		a.VideoSource = other.VideoSource
	}
	if !a.VideoCodec.Valid {
		a.VideoCodec = other.VideoCodec
	}
	if !a.Video3d.Valid {
		a.Video3d = other.Video3d
	}
	if !a.VideoModifier.Valid {
		a.VideoModifier = other.VideoModifier
	}
	if !a.ReleaseGroup.Valid {
		a.ReleaseGroup = other.ReleaseGroup
	}
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

func (a *ContentAttributes) InferVideoAttributes(input string) {
	a.VideoResolution = model.InferVideoResolution(input)
	a.VideoSource = model.InferVideoSource(input)
	a.VideoCodec, a.ReleaseGroup = model.InferVideoCodecAndReleaseGroup(input)
	a.Video3d = model.InferVideo3d(input)
	a.VideoModifier = model.InferVideoModifier(input)
}
