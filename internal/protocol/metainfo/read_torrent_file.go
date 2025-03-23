package metainfo

import (
	"fmt"

	"github.com/anacrolix/torrent/bencode"
)

type TorrentFile struct {
	Info Info `bencode:"info"`

	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list,omitempty"`
	CreationDate int64       `bencode:"creation date,omitempty"`
	Comment      string      `bencode:"comment,omitempty"`
	CreatedBy    string      `bencode:"created by,omitempty"`
	URLList      interface{} `bencode:"url-list,omitempty"`
}

func ReadTorrentFileBytes(bytes []byte) (TorrentFile, error) {
	var torrentFile TorrentFile
	if unmarshalErr := bencode.Unmarshal(bytes, &torrentFile); unmarshalErr != nil {
		return TorrentFile{}, fmt.Errorf("error unmarshaling torrent file: %s", unmarshalErr)
	}
	return torrentFile, nil
}
