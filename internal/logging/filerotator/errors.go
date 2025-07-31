package filerotator

import "errors"

var (
	Err                = errors.New("log file rotator")
	ErrStart           = errors.New("failed to start")
	ErrShutdown        = errors.New("failed to shutdown")
	ErrWriteFile       = errors.New("failed to write file")
	ErrFlush           = errors.New("failed to flush file")
	ErrCloseFile       = errors.New("failed to close file")
	ErrCreateFile      = errors.New("failed to create file")
	ErrCreateDirectory = errors.New("failed to create directory")
	ErrRotate          = errors.New("failed to rotate")
	ErrPrune           = errors.New("failed to prune")
	ErrReadDirectory   = errors.New("failed to read directory")
)
