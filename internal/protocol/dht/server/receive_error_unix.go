//go:build !windows

package server

func ignoreReadFromError(error) bool {
	// Good unix.
	return false
}
