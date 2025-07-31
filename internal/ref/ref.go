package ref

import (
	"errors"
	"regexp"
	"strings"
	"sync"
)

var mtx sync.Mutex

type Ref struct {
	name   string
	valid  bool
	parent *Ref
	subs   map[string]Ref
}

var Root = Ref{
	name:  "[root]",
	valid: true,
	subs:  make(map[string]Ref),
}

func NewNameNode(name string) Ref {
	return Ref{
		name: name,
	}
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
	if len(parts) == 0 {
		return Ref{}, errors.New("invalid name")
	}

	result := NewNameNode(parts[0])

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

func (n Ref) Sub(name string) (Ref, error) {
	mtx.Lock()
	defer mtx.Unlock()

	return n.sub(name)
}

func (n Ref) sub(name string) (Ref, error) {
	if !regexName.MatchString(name) {
		return Ref{}, errors.New("invalid name")
	}

	if _, ok := n.subs[name]; ok {
		return Ref{}, errors.New("name already exists")
	}

	return Ref{
		name:   name,
		parent: &n,
		valid:  n.valid,
	}, nil
}

func (n Ref) MustSub(name string) Ref {
	sub, err := n.Sub(name)
	if err != nil {
		panic(err)
	}

	return sub
}

func (n Ref) Name() string {
	return n.name
}

func (n Ref) Path() []string {
	var result []string
	for current := &n; current != nil && !current.IsRoot(); current = current.parent {
		result = append([]string{current.name}, result...)
	}

	return result
}

func (n Ref) IsRoot() bool {
	return n.parent == nil
}

func (n Ref) IsDescendentOf(other Ref) bool {
	for current := n.parent; current != nil; current = current.parent {
		if current.String() == other.String() {
			return true
		}
	}

	return false
}

func (n Ref) IsValid() bool {
	return n.valid
}

func (n Ref) String() string {
	return strings.Join(n.Path(), ".")
}

func (r *Ref) UnmarshalText(text []byte) error {
	ref, err := Parse(string(text))
	if err != nil {
		return err
	}
	*r = ref

	return nil
}
