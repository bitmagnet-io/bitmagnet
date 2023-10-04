package dht

import (
	"github.com/bitmagnet-io/bitmagnet/internal/dht/routing"
	"time"
)

type Config struct {
	PeerID                      [20]byte
	Routing                     routing.Config
	CrawlBootstrapHostsInterval time.Duration
	SampleInfoHashesInterval    time.Duration
	// DiscardUnscrapableTorrents when true, torrents that cannot be scraped to find seeders and leechers will be discarded
	DiscardUnscrapableTorrents bool
	MaxStagingSize             uint
	// SavePieces when true, torrent pieces will be persisted to the database.
	// The pieces take up quite a lot of space, and aren't currently very useful, but they may be used by future features.
	SavePieces bool
	// RescrapeThreshold is the amount of time that must pass before a torrent is rescraped to count seeders and leechers.
	RescrapeThreshold time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		PeerID:                      RandomPeerID(),
		Routing:                     routing.NewDefaultConfig(),
		CrawlBootstrapHostsInterval: time.Minute,
		SampleInfoHashesInterval:    time.Millisecond * 100,
		DiscardUnscrapableTorrents:  false,
		MaxStagingSize:              250,
		SavePieces:                  false,
		RescrapeThreshold:           time.Hour * 24 * 7,
	}
}
