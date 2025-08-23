package file_rotator

import (
	"bufio"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/fs"
)

type fileRotatorFile struct {
	writer *bufio.Writer
	file   fs.File
}

func (f *fileRotatorFile) write(p []byte) (int, error) {
	n, err := f.writer.Write(p)
	if err != nil {
		return n, fmt.Errorf("%w: %w", ErrWriteFile, err)
	}

	return n, nil
}

func (f *fileRotatorFile) flush() error {
	err := f.writer.Flush()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", ErrWriteFile, ErrFlush, err)
	}

	return err
}

func (f *fileRotatorFile) close() error {
	err := f.flush()

	if err == nil {
		err = f.file.Close()
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrCloseFile, err)
	}

	return nil
}
