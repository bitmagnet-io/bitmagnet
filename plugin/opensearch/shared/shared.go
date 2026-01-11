//go:build wasip1

package shared

import "errors"

type RecordType string

const (
	RecordTorrent RecordType = "torrents"
	RecordContent RecordType = "content"
)

var Err = errors.New("plugin-opensearch")
