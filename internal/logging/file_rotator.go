package logging

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func newFileRotator(
	config FileRotatorConfig,
) *fileRotator {
	return &fileRotator{
		path:       config.Path,
		baseName:   config.BaseName,
		maxAge:     config.MaxAge,
		maxSize:    config.MaxSize,
		maxBackups: config.MaxBackups,
		bufferSize: config.BufferSize,
	}
}

type fileRotator struct {
	lock        sync.Mutex
	path        string
	pathCreated bool
	baseName    string
	maxAge      time.Duration
	maxSize     int
	maxBackups  int
	bufferSize  int
	size        int
	nextTime    time.Time
	file        *fileRotatorFile
	closed      bool
}

func (r *fileRotator) Write(output []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.closed {
		return len(output), nil
	}

	if !r.pathCreated {
		err := os.MkdirAll(r.path, 0o755)
		if err != nil {
			return 0, err
		}

		r.pathCreated = true
	}

	if err := r.checkRotate(len(output)); err != nil {
		return 0, err
	}

	n, err := r.file.Write(output)
	if err != nil {
		return 0, err
	}

	r.size += n

	return n, nil
}

func (r *fileRotator) Sync() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.file == nil {
		return nil
	}

	return r.file.Close()
}

func (r *fileRotator) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.closed = true
	if r.file == nil {
		return nil
	}

	return r.file.Close()
}

func (r *fileRotator) checkRotate(n int) error {
	if !r.shouldRotate(n) {
		return nil
	}

	return r.rotate()
}

func (r *fileRotator) shouldRotate(n int) bool {
	if r.file == nil {
		return true
	}

	if r.maxAge > 0 && time.Now().After(r.nextTime) {
		return true
	}

	if r.maxSize > 0 && r.size+n > r.maxSize {
		return true
	}

	return false
}

func (r *fileRotator) rotate() error {
	if r.file != nil {
		err := r.file.Close()
		r.file = nil

		if err != nil {
			return err
		}
	}

	now := time.Now()

	fp, err := newFileRotatorFile(r.newFilePath(now), r.bufferSize)
	if err != nil {
		return err
	}

	r.file = fp
	r.size = 0
	r.nextTime = now.Add(r.maxAge)

	return r.pruneBackups(now)
}

const timeFormat = "2006-01-02-15-04-05"

func (r *fileRotator) newFilePath(now time.Time) string {
	return path.Join(r.path, fmt.Sprintf("%s.%s.log", r.baseName, now.Format(timeFormat)))
}

func (r *fileRotator) pruneBackups(now time.Time) error {
	files, err := os.ReadDir(r.path)
	if err != nil {
		return err
	}

	//nolint:prealloc
	var backupFiles []string

	strNow := now.Format(timeFormat)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasPrefix(name, r.baseName+".") || !strings.HasSuffix(name, ".log") {
			continue
		}

		strDate := name[len(r.baseName)+1 : len(name)-4]

		_, parseErr := time.Parse(timeFormat, strDate)
		if parseErr != nil {
			continue
		}

		if strDate >= strNow {
			continue
		}

		backupFiles = append(backupFiles, name)
	}

	if len(backupFiles) <= r.maxBackups {
		return nil
	}

	for _, name := range backupFiles[:len(backupFiles)-r.maxBackups] {
		// make a best effort to remove the file
		_ = os.Remove(path.Join(r.path, name))
	}

	return nil
}

func newFileRotatorFile(path string, bufferSize int) (*fileRotatorFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
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

func (f *fileRotatorFile) Write(p []byte) (int, error) {
	return f.writer.Write(p)
}

func (f *fileRotatorFile) Flush() error {
	return f.writer.Flush()
}

func (f *fileRotatorFile) Close() error {
	if err := f.Flush(); err != nil {
		return err
	}

	return f.file.Close()
}
