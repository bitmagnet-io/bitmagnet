package target

import "errors"

var (
	Err               = errors.New("target")
	ErrUnknownTarget  = errors.New("unknown target")
	ErrLookupTorrents = errors.New("lookup torrents failed")
	ErrNoTorrents     = errors.New("no torrents found")
	ErrDataSchema     = errors.New("data schema failed")
	ErrUISchema       = errors.New("uischema failed")
	ErrMarshalData    = errors.New("marshal data failed")
	ErrValidation     = errors.New("data validation failed")
	ErrSend           = errors.New("send failed")
	ErrPlugin         = errors.New("plugin returned an error")
)
