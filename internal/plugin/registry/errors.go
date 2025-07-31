package registry

import "errors"

var (
	Err              = errors.New("registry")
	ErrResolve       = errors.New("resolve failed")
	ErrUnknownPlugin = errors.New("unknown plugin")
	ErrDisabled      = errors.New("plugin is disabled")
	ErrDependency    = errors.New("plugin dependency")
)
