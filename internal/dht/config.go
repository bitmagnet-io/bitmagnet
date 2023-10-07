package dht

import (
	"time"
)

type Config struct {
	Routing                     RoutingTableConfig
	CrawlBootstrapHostsInterval time.Duration
	SampleInfoHashesInterval    time.Duration
	// DiscardUnscrapableTorrents when true, torrents that cannot be scraped to find seeders and leechers will be discarded
	DiscardUnscrapableTorrents bool
	MaxStagingSize             uint
	// SaveFiles when true, torrent files metadata will be persisted to the database.
	SaveFiles bool
	// SaveFilesThreshold specifies a maximum number of files in a torrent before file information is discarded.
	// Some torrents contain thousands of files which can severely impact performance and uses a lot of disk space.
	// A value of 0 means no threshold.
	SaveFilesThreshold uint
	// SavePieces when true, torrent pieces will be persisted to the database.
	// The pieces take up quite a lot of space, and aren't currently very useful, but they may be used by future features.
	SavePieces bool
	// RescrapeThreshold is the amount of time that must pass before a torrent is rescraped to count seeders and leechers.
	RescrapeThreshold time.Duration
}

type RoutingTableConfig struct {
	MaxPeers              uint
	MaxConcurrency        uint
	MaxConcurrencyPerPeer uint
}

func NewDefaultConfig() Config {
	return Config{
		Routing: RoutingTableConfig{
			MaxPeers:              1000,
			MaxConcurrency:        500,
			MaxConcurrencyPerPeer: 2,
		},
		CrawlBootstrapHostsInterval: time.Minute,
		SampleInfoHashesInterval:    time.Millisecond * 100,
		DiscardUnscrapableTorrents:  false,
		MaxStagingSize:              250,
		SaveFiles:                   true,
		SaveFilesThreshold:          50,
		SavePieces:                  false,
		RescrapeThreshold:           time.Hour * 24 * 7,
	}
}
