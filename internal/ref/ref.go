package ref

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

var mtx sync.Mutex

type Ref struct {
	name      string
	str       string
	canonical bool
	parent    *Ref
	subs      map[string]Ref
}

var Root = Ref{
	name:      "[root]",
	canonical: true,
	subs:      make(map[string]Ref),
}

func MustNew(name string) Ref {
	ref, err := New(name)
	if err != nil {
		panic(err)
	}

	return ref
}

func New(name string) (Ref, error) {
	if !regexName.MatchString(name) {
		return Ref{}, fmt.Errorf("%w: %w: %s", Err, ErrInvalidName, name)
	}

	return Ref{
		name: name,
		subs: make(map[string]Ref),
	}, nil
}

func MustParse(name string) Ref {
	ref, err := Parse(name)
	if err != nil {
		panic(err)
	}

	return ref
}

func Parse(name string) (Ref, error) {
	parts := strings.Split(name, ".")

	result, err := New(parts[0])
	if err != nil {
		return Ref{}, err
	}

	for i := 0; i < len(parts); i++ {
		nextResult, err := result.sub(parts[i])
		if err != nil {
			return Ref{}, err
		}
		result = nextResult
	}

	return result, nil
}

var regexName = regexp.MustCompile(`^[a-z0-9]+(?:_[a-z0-9]+)*$`)

func (r Ref) Sub(name string) (Ref, error) {
	mtx.Lock()
	defer mtx.Unlock()

	return r.sub(name)
}

func (r Ref) sub(name string) (Ref, error) {
	if !regexName.MatchString(name) {
		return Ref{}, fmt.Errorf("%w: %w: %s.%s", Err, ErrInvalidName, r.String(), name)
	}

	if _, ok := r.subs[name]; ok {
		return Ref{}, fmt.Errorf("%w: %w: %s.%s", Err, ErrNameAlreadyExists, r.String(), name)
	}

	str := r.str
	if !r.IsRoot() {
		str += "."
	}
	str += name

	sub := Ref{
		name:      name,
		str:       str,
		parent:    &r,
		subs:      make(map[string]Ref),
		canonical: r.canonical,
	}

	r.subs[name] = sub

	return sub, nil
}

func (r Ref) MustSub(name string) Ref {
	sub, err := r.Sub(name)
	if err != nil {
		panic(err)
	}

	return sub
}

func (r Ref) GetSub(name string) (Ref, bool) {
	sub, ok := r.subs[name]

	return sub, ok
}

func (r Ref) Name() string {
	return r.name
}

func (r Ref) Path() []string {
	var result []string
	for current := &r; !current.IsRoot(); current = current.parent {
		result = append([]string{current.name}, result...)
	}

	return result
}

func (r Ref) IsRoot() bool {
	return r.parent == nil
}

func (r Ref) IsDescendentOf(other Ref) bool {
	for current := r.parent; current != nil; current = current.parent {
		if current.String() == other.String() {
			return true
		}
	}

	return false
}

func (r Ref) IsCanonical() bool {
	return r.canonical
}

func (r Ref) String() string {
	return r.str
}

func (r Ref) Equals(other Ref) bool {
	return r.str == other.str
}

func (r *Ref) UnmarshalText(text []byte) error {
	ref, err := Parse(string(text))
	if err != nil {
		return err
	}
	*r = ref

	return nil
}

func (r Ref) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprintf(w, `"%s"`, r.String())
}

func (r *Ref) UnmarshalGQL(v any) error {
	vStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", v)
	}

	parsed, err := Parse(vStr)
	if err != nil {
		return err
	}

	*r = parsed

	return nil
}
