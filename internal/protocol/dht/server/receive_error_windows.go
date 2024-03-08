package server

import (
	"errors"
	"syscall"
)

func ignoreReadFromError(err error) bool {
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		return false
	}
	return errno == syscall.WSAECONNRESET
}
