package configresolver

import "reflect"

type Resolver interface {
	Key() string
	Priority() int
	Resolve(path []string, valueType reflect.Type) (interface{}, bool, error)
}

type baseResolver struct {
	key      string
	priority int
}

func (r *baseResolver) Key() string {
	return r.key
}

func (r *baseResolver) Priority() int {
	return r.priority
}

type Option struct {
	apply func(*baseResolver)
}

func WithKey(key string) Option {
	return Option{
		apply: func(r *baseResolver) {
			r.key = key
		},
	}
}

func WithPriority(priority int) Option {
	return Option{
		apply: func(r *baseResolver) {
			r.priority = priority
		},
	}
}

func (r *baseResolver) applyOptions(options ...Option) {
	for _, o := range options {
		o.apply(r)
	}
}
