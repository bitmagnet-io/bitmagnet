package model

type InfoHashStringer interface {
	InfoHashString() string
}

func (t Torrent) InfoHashString() string {
	return t.InfoHash.String()
}

func (t TorrentFile) InfoHashString() string {
	return t.InfoHash.String()
}

func (tc TorrentContent) InfoHashString() string {
	return tc.InfoHash.String()
}
