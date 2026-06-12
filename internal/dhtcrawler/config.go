package dhtcrawler

import (
	"time"
)

type Config struct {
	// ScalingFactor is a rough proxy for resource usage of the crawler; concurrency and buffer size of the various
	// pipeline channels are multiplied by this value. Diminishing returns may result from exceeding the
	// default value of 10. Since the software has not been tested on a wide variety of hardware and network
	// conditions; your mileage may vary here...
	ScalingFactor                uint
	BootstrapNodes               []string
	ReseedBootstrapNodesInterval time.Duration
	// SaveFilesThreshold specifies a maximum number of files in a torrent before file information is discarded.
	// Some torrents contain thousands of files which can severely impact performance and uses a lot of disk space.
	SaveFilesThreshold uint
	// SavePieces when true, torrent pieces will be persisted to the database.
	// The pieces take up quite a lot of space, and aren't currently very useful,
	// but they may be used by future features.
	SavePieces bool
	// RescrapeThreshold is the amount of time that must pass before a torrent is rescraped
	// to count seeders and leechers.
	RescrapeThreshold time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		ScalingFactor:                10,
		BootstrapNodes:               defaultBootstrapNodes,
		ReseedBootstrapNodesInterval: time.Minute,
		SaveFilesThreshold:           100,
		SavePieces:                   false,
		RescrapeThreshold:            time.Hour * 24 * 30,
	}
}

// https://github.com/anacrolix/dht/blob/92b36a3fa7a37a15e08684337b47d8d0fb322ab6/dht.go#L106
var defaultBootstrapNodes = []string{
    "dht.transmissionbt.com:6881"
    "router.bittorrent.com:6881"
    "router.utorrent.com:6881"
    "dht.vuze.com:6881"
    "router.bt.ouinet.work:6881"
    "router.bittorrent.com:8991"
    "router.bitcomet.com:688"
    "dht.aelitis.com:6881"
    "node1.hyperdht.org:49737"
    "node2.hyperdht.org:49737"
    "node3.hyperdht.org:49737"
    "dht.anacrolix.link:42069"
    "router.bittorrent.cloud:42069"
    "router.silotis.us:6881"
    "ntp.juliusbeckmann.de:6881"
    "mgts.ivth.ru:57858"
    "sorcerer.leentje.org:49786"
    "libertalia.space:50005"
    "milda.intelib.org:51413"
    "routerx.bt.ouinet.work:5060"
    "bootstrap.jami.net:4222"
    "dht.libtorrent.org:25401"
}
