package server

import (
	"errors"
	"golang.org/x/sys/windows"
	"syscall"
)

func ignoreReadFromError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case
			windows.WSAENETRESET,
			windows.WSAECONNRESET,
			windows.WSAECONNABORTED,
			windows.WSAECONNREFUSED,
			windows.WSAENETUNREACH,
			windows.WSAETIMEDOUT:
			return true
		}
	}
	return false
}
