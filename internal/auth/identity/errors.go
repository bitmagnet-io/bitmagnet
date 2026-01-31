package identity

import "errors"

var (
	Err               = errors.New("identity provider")
	ErrAuthentication = errors.New("authentication failed")
	ErrUnmatched      = errors.New("no provider matched")
)
