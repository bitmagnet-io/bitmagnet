package filerotator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func newFileRotatorFile(path string, bufferSize int) (*fileRotatorFile, error) {
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateDirectory, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateFile, err)
	}

	return &fileRotatorFile{
		writer: bufio.NewWriterSize(f, bufferSize),
		file:   f,
	}, nil
}

type fileRotatorFile struct {
	writer *bufio.Writer
	file   *os.File
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
