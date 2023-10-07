package metainfo

import (
	"errors"
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/anacrolix/torrent/bencode"
	mi "github.com/anacrolix/torrent/metainfo"
)

func ParseMetaInfoBytes(infoHash krpc.ID, metaInfoBytes []byte) (info mi.Info, err error) {
	if krpc.ID(mi.HashBytes(metaInfoBytes)) != infoHash {
		err = errors.New("info bytes have wrong hash")
		return
	}
	if unmarshalErr := bencode.Unmarshal(metaInfoBytes, &info); unmarshalErr != nil {
		err = fmt.Errorf("error unmarshaling info bytes: %s", unmarshalErr)
		return
	}
	return
}
