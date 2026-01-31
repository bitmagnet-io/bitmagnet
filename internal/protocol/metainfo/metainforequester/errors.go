package metainforequester

import "errors"

var (
	Err                 = errors.New("meta info requester")
	ErrAcquireSemaphore = errors.New("failed to acquire semaphore")
	ErrConnect          = errors.New("connection failed")
	ErrDial             = errors.New("dial failed")
	ErrSetLinger        = errors.New("set linger failed")
	ErrSetDeadline      = errors.New("set deadline failed")
	ErrHandshake        = errors.New("handshake failed")
	ErrWrite            = errors.New("write failed")
	ErrRead             = errors.New("read failed")
	ErrUnsupported      = errors.New("unsupported")
	ErrRejected         = errors.New("rejected")
	ErrInvalidResponse  = errors.New("invalid response")
	ErrFirstExMessage   = errors.New("first extension message is not an extension handshake")
	ErrUnmarshal        = errors.New("unmarshal failed")
	ErrSize             = errors.New("invalidd metadata size")
	ErrUTMetadata       = errors.New("ut_metadata is not an uint8")
	ErrInfoHashMismatch = errors.New("info hash mismatch")
	ErrExHandshake      = errors.New("extended handshake failed")
	ErrRequestPieces    = errors.New("pieces request failed")
	ErrReadPieces       = errors.New("failed to read pieces")
	ErrParse            = errors.New("failed to parse meta info")
)
