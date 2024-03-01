package dhtcrawler

import (
	"github.com/anacrolix/torrent/bencode"
	"github.com/adrg/xdg"
	"os"
	"path"
)

func (c *crawler) saveTorrentsInfo(torrentsInfo []infoHashWithMetaInfo) {
	if (c.saveFailed) {
		for _, info := range torrentsInfo {
			torrentsFolder := path.Join(xdg.StateHome, "bitmagnet", "torrents")
			err := os.MkdirAll(torrentsFolder, 0755)
			if err != nil {
				c.logger.Errorf("failed to create folder for torrents! Error %s", err)
				return
			}

			torrentFile := path.Join(torrentsFolder, info.infoHash.String() + ".torrent")
			io, err := os.Create(torrentFile)
			if err != nil {
				c.logger.Errorf("failed to create torrent file! Error %s", err)
				return
			}

			encoder := bencode.NewEncoder(io)
			encoder.Encode(info.metaInfo)

			if err := io.Close(); err != nil {
				c.logger.Errorf("failed to close file! Error %s", err)
			}
			c.logger.Errorf("saved %s", torrentFile)
		}
	}
}
