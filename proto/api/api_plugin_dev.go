//go:build !wasip1

package api

func RegisterHTTPHandler(HTTPHandler) {
	panic("not implemented")
}

func RegisterPlugin(Plugin) {
	panic("not implemented")
}

func RegisterTorrentTarget(TorrentTarget) {
	panic("not implemented")
}
