package dhtcrawler

import (
	"time"
)

type Config struct {
	ScalingFactor               uint
	CrawlBootstrapHostsInterval time.Duration
	// SaveFilesThreshold specifies a maximum number of files in a torrent before file information is discarded.
	// Some torrents contain thousands of files which can severely impact performance and uses a lot of disk space.
	SaveFilesThreshold uint
	// SavePieces when true, torrent pieces will be persisted to the database.
	// The pieces take up quite a lot of space, and aren't currently very useful, but they may be used by future features.
	SavePieces bool
	// RescrapeThreshold is the amount of time that must pass before a torrent is rescraped to count seeders and leechers.
	RescrapeThreshold time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		ScalingFactor:               10,
		CrawlBootstrapHostsInterval: time.Minute,
		SaveFilesThreshold:          50,
		SavePieces:                  false,
		RescrapeThreshold:           time.Hour * 24 * 30,
	}
}
