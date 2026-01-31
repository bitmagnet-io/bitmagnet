package error_registry

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type Registry struct {
	refMap ref.Map[error]
}

func New(options ...Option) Registry {
	r := &Registry{
		refMap: ref.NewMap[error](),
	}

	for _, option := range options {
		option(r)
	}

	return *r
}

type Option func(*Registry)

func WithEntry(ref ref.Ref, err error) Option {
	return func(r *Registry) {
		r.refMap.Set(ref, err)
	}
}

func WithEntries(entries ref.Map[error]) Option {
	return func(r *Registry) {
		r.refMap.SetAll(entries)
	}
}

func (r Registry) Get(ref ref.Ref) (ref.Entry[error], bool) {
	return r.refMap.EntryOK(ref)
}

func (r Registry) Entries() []ref.Entry[error] {
	return r.refMap.Entries()
}

func (r Registry) Identify(err error) (ref.Entry[error], bool) {
	for _, entry := range r.Entries() {
		if errors.Is(err, entry.Value) {
			return entry, true
		}
	}

	return ref.Entry[error]{}, false
}
