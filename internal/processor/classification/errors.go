package classification

import (
	"fmt"
	"strings"
)

type Error interface {
	error
	Key() string
}

type WorkflowError struct {
	key     string
	message string
	cause   error
}

func (e WorkflowError) Error() string {
	if e.message != "" {
		return e.message
	}
	return fmt.Sprintf("workflow unmarshalError: %s", e.key)
}

func (e WorkflowError) Key() string {
	return e.key
}

var ErrNoMatch = WorkflowError{
	key: "no_match",
}

var ErrDeleteTorrent = WorkflowError{
	key: "delete_torrent",
}

type RuntimeError struct {
	Path  []string
	Cause error
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at Path %s: %s", strings.Join(e.Path, "."), e.Cause)
}

func (e RuntimeError) Unwrap() error {
	return e.Cause
}
