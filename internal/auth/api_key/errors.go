package api_key

import "errors"

var (
	Err           = errors.New("api key")
	ErrCreate     = errors.New("create failed")
	ErrAuth       = errors.New("auth failed")
	ErrDecode     = errors.New("decode failed")
	ErrMismatch   = errors.New("mismatch")
	ErrLength     = errors.New("wrong length")
	ErrList       = errors.New("list failed")
	ErrDelete     = errors.New("delete failed")
	ErrNotFound   = errors.New("not found")
	ErrExpired    = errors.New("expired")
	ErrRepository = errors.New("repository")
)
