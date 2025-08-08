package ref

import "errors"

var (
	Err                  = errors.New("ref")
	ErrInvalidName       = errors.New("invalid name")
	ErrNameAlreadyExists = errors.New("name already exists")
)
