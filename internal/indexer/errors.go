package indexer

import (
	"errors"
)

var (
	Err                  = errors.New("indexer")
	ErrSetup             = errors.New("setup")
	ErrInterrupted       = errors.New("interrupted")
	ErrClassify          = errors.New("classify")
	ErrPersist           = errors.New("persist")
	ErrShutdown          = errors.New("shutdown")
	ErrAllTorrentsFailed = errors.New("all torrents failed")
)
