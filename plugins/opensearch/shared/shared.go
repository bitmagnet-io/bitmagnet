package shared

import "errors"

type RecordType string

const (
	TemplateVersion = 1

	RecordTorrent RecordType = "torrents"
	RecordContent RecordType = "content"
)

var Err = errors.New("opensearch")
