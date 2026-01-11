package search

import "context"

type Base interface {
	base()
}

type TorrentContent interface {
	TorrentContent(ctx context.Context, params Params) (TorrentContentResult, error)
}

type TorrentFiles interface {
	TorrentFiles(ctx context.Context, params Params) (TorrentFilesResult, error)
}

type Content interface {
	Content(ctx context.Context, params Params) (ContentResult, error)
}

type Search interface {
	Base
	TorrentContent
	TorrentFiles
	Content
}
