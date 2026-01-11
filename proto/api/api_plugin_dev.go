//go:build !wasip1

package api

func RegisterHTTPHandler(HTTPHandler) {
	panic("not implemented")
}
