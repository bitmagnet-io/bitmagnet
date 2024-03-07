package workflow

import (
	"fmt"
	"strings"
)

type Error interface {
	error
	Key() string
}

type workflowError struct {
	key     string
	message string
	cause   error
}

func (e workflowError) Error() string {
	if e.message != "" {
		return e.message
	}
	return fmt.Sprintf("workflow unmarshalError: %s", e.key)
}

func (e workflowError) Key() string {
	return e.key
}

var ErrNoMatch = workflowError{
	key: "no_match",
}

var ErrDeleteTorrent = workflowError{
	key: "delete_torrent",
}

type runtimeError struct {
	path  []string
	cause error
}

func (e runtimeError) Error() string {
	return fmt.Sprintf("runtime error at path %s: %s", strings.Join(e.path, "."), e.cause)
}

func (e runtimeError) Unwrap() error {
	return e.cause
}
