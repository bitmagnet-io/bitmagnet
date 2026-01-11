package env

import (
	"context"
	"io"
	"os"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
)

type (
	Context context.Context
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer

	Reader interface {
		Stdin
	}

	Writer interface {
		Stdout
		Stderr() Stderr
	}

	ReadWriter interface {
		Reader
		Writer
	}

	VarsLookup interface {
		LookupVar(key string) (string, bool)
	}

	VarsProvider func() map[string]string

	ArgsProvider interface {
		Args() []string
	}

	SignalsProvider interface {
		Signals(...os.Signal) <-chan os.Signal
	}

	Env interface {
		Context
		ReadWriter
		VarsLookup
		ArgsProvider
		SignalsProvider
		fs.FSProvider
	}
)

func New(options ...Option) Env {
	e := environment{
		env: make(map[string]string),
	}

	Options(append(options, ensureValues())...)(&e)

	return e
}

func NewDefault() Env {
	return New(WithDefaults())
}

type environment struct {
	context.Context
	io.Reader
	io.Writer
	SignalsProvider
	fs.FSProvider
	env    map[string]string
	args   []string
	stderr io.Writer
}

func (e environment) Args() []string {
	return slices.Collect(slices.Values(e.args))
}

func (e environment) Stderr() Stderr {
	return e.stderr
}

func (e environment) LookupVar(key string) (string, bool) {
	value, ok := e.env[key]

	return value, ok
}
