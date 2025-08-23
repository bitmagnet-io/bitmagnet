package lookup

import (
	"errors"
)

var (
	Err                 = errors.New("config lookup")
	ErrReadFile         = errors.New("failed to read file")
	ErrYAMLUnmarshal    = errors.New("failed to unmarshal YAML")
	ErrLookup           = errors.New("lookup failed")
	ErrInvalidStructure = errors.New("invalid config structure")
)
