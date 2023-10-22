package metainfo

import (
	"errors"
	"fmt"
	"github.com/anacrolix/torrent/bencode"
	mi "github.com/anacrolix/torrent/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"unicode/utf8"
)

func ParseMetaInfoBytes(infoHash protocol.ID, metaInfoBytes []byte) (Info, error) {
	if protocol.ID(mi.HashBytes(metaInfoBytes)) != infoHash {
		return Info{}, errors.New("info bytes have wrong hash")
	}
	var info Info
	if unmarshalErr := bencode.Unmarshal(metaInfoBytes, &info); unmarshalErr != nil {
		return Info{}, fmt.Errorf("error unmarshaling info bytes: %s", unmarshalErr)
	}
	checkUtf8Strings := make([]string, 0, len(info.Files)+1)
	checkUtf8Strings = append(checkUtf8Strings, info.BestName())
	for _, file := range info.Files {
		checkUtf8Strings = append(checkUtf8Strings, file.DisplayPath(&info))
	}
	for _, str := range checkUtf8Strings {
		if !utf8.ValidString(str) {
			return Info{}, errors.New("meta info contains an invalid utf8 string")
		}
	}
	return info, nil
}
