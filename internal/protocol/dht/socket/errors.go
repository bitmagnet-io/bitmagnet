package socket

import (
	"errors"
)

var (
	Err                       = errors.New(Namespace)
	ErrUnknownAdapter         = errors.New("unknown adapter")
	ErrOpenFailed             = errors.New("open failed")
	ErrCreateFailed           = errors.New("create failed")
	ErrBindFailed             = errors.New("bind failed")
	ErrSetOptionFailed        = errors.New("set option failed")
	ErrCloseFailed            = errors.New("close failed")
	ErrClosed                 = errors.New("closed")
	ErrInvalidAddress         = errors.New("invalid address")
	ErrUnsupportedAddressType = errors.New("unsupported address type")
	ErrSendFailed             = errors.New("send failed")
	ErrReceiveFailed          = errors.New("receive failed")
)
