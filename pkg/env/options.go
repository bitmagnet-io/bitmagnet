package env

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
)

type Option func(*environment)

func Options(options ...Option) Option {
	return func(e *environment) {
		for _, option := range options {
			option(e)
		}
	}
}

func WithDefaults() Option {
	return Options(
		WithContext(context.Background()),
		WithOSReadWriters(),
		WithOSArgs(),
		WithOSSignalsProvider(),
		WithVarsFromEnviron(),
		WithVarsFromDotEnv(),
		WithXDGFSProvider("bitmagnet"),
	)
}

func WithContext(ctx context.Context) Option {
	return func(e *environment) {
		e.Context = ctx
	}
}

func WithArgs(args ...string) Option {
	return func(e *environment) {
		e.args = args
	}
}

func WithOSArgs() Option {
	return WithArgs(os.Args[1:]...)
}

func WithStdin(stdin io.Reader) Option {
	return func(e *environment) {
		e.Reader = stdin
	}
}

func WithStdout(stdout io.Writer) Option {
	return func(e *environment) {
		e.Writer = stdout
	}
}

func WithStderr(stderr io.Writer) Option {
	return func(e *environment) {
		e.stderr = stderr
	}
}

func WithOSReadWriters() Option {
	return Options(
		WithStdin(os.Stdin),
		WithStdout(os.Stdout),
		WithStderr(os.Stderr),
	)
}

func WithVarsProvider(provider VarsProvider) Option {
	return func(e *environment) {
		for k, v := range provider() {
			e.env[k] = v
		}
	}
}

func WithVarsFromEnviron() Option {
	return WithVarsProvider(providerOSEnviron())
}

func WithVarsFromDotEnv() Option {
	return WithVarsProvider(providerDotEnv())
}

func WithSignalsProvider(provider SignalsProvider) Option {
	return func(e *environment) {
		e.SignalsProvider = provider
	}
}

func WithOSSignalsProvider() Option {
	return func(e *environment) {
		e.SignalsProvider = osSignalsProvider{}
	}
}

func WithFSProvider(provider fs.Provider) Option {
	return func(e *environment) {
		e.Provider = provider
	}
}

func WithXDGFSProvider(subPath string) Option {
	return func(e *environment) {
		e.Provider = fs.NewFSProviderXDG(subPath)
	}
}

func ensureValues() Option {
	return func(e *environment) {
		if e.Context == nil {
			e.Context = context.TODO()
		}

		if e.Writer == nil {
			e.Writer = io.Discard
		}

		if e.Reader == nil {
			e.Reader = &bytes.Buffer{}
		}

		if e.stderr == nil {
			e.stderr = e.Writer
		}

		if e.SignalsProvider == nil {
			e.SignalsProvider = nopSignalsProvider{}
		}

		if e.Provider == nil {
			e.Provider = fs.ProviderNop{}
		}
	}
}
