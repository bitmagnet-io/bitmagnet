package httpserver

import (
	"cmp"

	"github.com/gin-gonic/gin"
)

type Phase int

const (
	PhasePre     = -1
	PhaseDefault = 0
	PhasePost    = 1
)

type Option interface {
	Key() string
	Phase() Phase
	Compare(other Option) int
	Apply(engine *gin.Engine)
}

func NewOption(key string, phase Phase, fn func(engine *gin.Engine)) Option {
	return option{
		key:   key,
		phase: phase,
		fn:    fn,
	}
}

type option struct {
	key   string
	phase Phase
	fn    func(engine *gin.Engine)
}

func (o option) Key() string {
	return o.key
}

func (o option) Phase() Phase {
	return o.phase
}

func (o option) Apply(engine *gin.Engine) {
	o.fn(engine)
}

func (o option) Compare(other Option) int {
	result := cmp.Compare(o.Phase(), other.Phase())
	if result != 0 {
		return result
	}

	return cmp.Compare(o.Key(), other.Key())
}
