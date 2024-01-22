package model

func (h TorrentHint) IsNil() bool {
	return h.ContentType.IsNil()
}

func (h TorrentHint) NullContentType() NullContentType {
	if h.IsNil() {
		return NullContentType{}
	}
	return NewNullContentType(h.ContentType)
}

func (h TorrentHint) ContentRef() Maybe[ContentRef] {
	if h.ContentID.Valid {
		return MaybeValid(ContentRef{
			Type:   h.ContentType,
			Source: h.ContentSource.String,
			ID:     h.ContentID.String,
		})
	}
	return Maybe[ContentRef]{}
}
