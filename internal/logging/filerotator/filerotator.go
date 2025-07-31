package filerotator

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func New(
	config Config,
) *FileRotator {
	return &FileRotator{
		path:       config.Path,
		baseName:   config.BaseName,
		maxAge:     config.MaxAge,
		maxSize:    config.MaxSize,
		maxBackups: config.MaxBackups,
		bufferSize: config.BufferSize,
	}
}

type FileRotator struct {
	mutex      sync.Mutex
	path       string
	baseName   string
	maxAge     time.Duration
	maxSize    int
	maxBackups int
	bufferSize int
	size       int
	nextTime   time.Time
	file       *fileRotatorFile
	started    bool
	err        chan error
}

func (r *FileRotator) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		if r.started {
			return runner.NopShutdowner, fmt.Errorf("%w: %w: %w", Err, ErrStart, runner.ErrAlreadyRunning)
		}

		r.started = true
		r.err = make(chan error, 1)

		shutdown := make(chan struct{})

		go func() {
			var err error

			select {
			case <-ctx.Done():
				err = fmt.Errorf("%w: %w", Err, ctx.Err())
			case <-shutdown:
			case err = <-r.err:
			}

			r.mutex.Lock()
			defer r.mutex.Unlock()

			r.started = false
			r.file = nil

			cancel(err)
		}()

		return func(context.Context) error {
			r.mutex.Lock()
			defer r.mutex.Unlock()

			if r.file != nil {
				err := r.file.close()
				if err != nil {
					return fmt.Errorf("%w: %w: %w", Err, ErrShutdown, err)
				}
			}

			r.started = false
			r.file = nil

			close(shutdown)

			return nil
		}, nil
	}
}

func (r *FileRotator) Write(output []byte) (int, error) {
	n, err := r.write(output)
	if err != nil {
		r.err <- err
	}

	return n, err
}

func (r *FileRotator) Sync() error {
	err := r.sync()
	if err != nil {
		r.err <- err
	}

	return err
}

func (r *FileRotator) write(output []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.started {
		return len(output), nil
	}

	if err := r.checkRotate(len(output)); err != nil {
		return 0, err
	}

	n, err := r.file.write(output)
	if err != nil {
		return 0, err
	}

	r.size += n

	return n, nil
}

func (r *FileRotator) checkRotate(n int) error {
	if !r.shouldRotate(n) {
		return nil
	}

	return r.rotate()
}

func (r *FileRotator) shouldRotate(n int) bool {
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

func (r *FileRotator) rotate() error {
	if r.file != nil {
		err := r.file.close()
		r.file = nil

		if err != nil {
			return fmt.Errorf("%w: %w: %w", Err, ErrRotate, err)
		}
	}

	now := time.Now()

	fp, err := newFileRotatorFile(r.newFilePath(now), r.bufferSize)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrRotate, err)
	}

	r.file = fp
	r.size = 0
	r.nextTime = now.Add(r.maxAge)

	return r.prune(now)
}

func (r *FileRotator) sync() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.file == nil {
		return nil
	}

	err := r.file.flush()
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	return nil
}

const timeFormat = "2006-01-02-15-04-05"

func (r *FileRotator) newFilePath(now time.Time) string {
	return path.Join(r.path, fmt.Sprintf("%s.%s.log", r.baseName, now.Format(timeFormat)))
}

func (r *FileRotator) prune(now time.Time) error {
	files, err := os.ReadDir(r.path)
	if err != nil {
		return fmt.Errorf("%w: %w: %w: %w", Err, ErrPrune, ErrReadDirectory, err)
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
