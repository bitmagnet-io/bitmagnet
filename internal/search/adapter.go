package search

import "context"

type TorrentContentAdapter interface {
	TorrentContent(ctx context.Context, params Params) (TorrentContentResult, error)
}

type TorrentsAdapter interface {
	Torrents(ctx context.Context, params Params) (TorrentsResult, error)
}

type TorrentFilesAdapter interface {
	TorrentFiles(ctx context.Context, params Params) (TorrentFilesResult, error)
}

type ContentAdapter interface {
	Content(ctx context.Context, params Params) (ContentResult, error)
}

type Adapter interface {
	TorrentContentAdapter
	TorrentFilesAdapter
	ContentAdapter
}
