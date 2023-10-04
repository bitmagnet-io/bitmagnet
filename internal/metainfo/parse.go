package metainfo

import (
	"errors"
	"fmt"
	"github.com/anacrolix/torrent/bencode"
	mi "github.com/anacrolix/torrent/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func ParseMetaInfoBytes(infoHash model.Hash20, metaInfoBytes []byte) (info mi.Info, err error) {
	if model.Hash20(mi.HashBytes(metaInfoBytes)) != infoHash {
		err = errors.New("info bytes have wrong hash")
		return
	}
	if unmarshalErr := bencode.Unmarshal(metaInfoBytes, &info); unmarshalErr != nil {
		err = fmt.Errorf("error unmarshaling info bytes: %s", unmarshalErr)
		return
	}
	return
}
