package resolver

import (
	"reflect"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"gopkg.in/yaml.v3"
)

type (
	Param struct {
		mtx sync.RWMutex
		param
		*chain
		pending bool
	}

	param = registry.Param

	chain struct {
		source       string
		value        any
		initialValue any
		prev         *chain
	}
)

func (p *Param) Source() string {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.source
}

func (p *Param) Value() any {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.value
}

func (p *Param) ValueString() string {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	str, err := p.StringifyAny(p.value)
	if err != nil {
		panic(err)
	}

	return str
}

func (p *Param) ValueYAML() yaml.Node {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	node, err := p.EncodeYAMLAny(p.value)
	if err != nil {
		panic(err)
	}

	return node
}

func (p *Param) Prev() (*Param, bool) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.prev == nil {
		return nil, false
	}

	return &Param{
		param: p.param,
		chain: p.prev,
	}, true
}

func (p *Param) Last() *Param {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	current := p.chain

	for {
		if current.prev == nil {
			return &Param{
				param: p.param,
				chain: current,
			}
		}

		current = current.prev
	}
}

func (p *Param) Save(value any) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	chainCopy := p.chain

	if reflected, ok := atomic.ReflectValue(reflect.ValueOf(p.value)); ok {
		reflected.SetAny(value)

		if p.chain.source != config.SourceDynamic {
			p.chain = &chain{
				source: config.SourceDynamic,
				value:  chainCopy.value,
				prev:   chainCopy,
			}
		}
	} else {
		if p.chain.source != config.SourcePending {
			p.chain = &chain{
				source: config.SourcePending,
				prev:   chainCopy,
			}
		}

		p.value = value
		p.pending = true
	}
}

func (p *Param) Delete() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if reflected, ok := atomic.ReflectValue(reflect.ValueOf(p.value)); ok {
		if p.source == config.SourceDynamic {
			reflected.SetAny(p.initialValue)
			p.chain = p.chain.prev
		}
	} else {
		if p.source == config.SourcePending || p.source == config.SourcePersisted {
			p.pending = p.source == config.SourcePersisted
			p.chain = p.chain.prev
		}
	}
}

func (p *Param) IsPending() bool {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.pending
}
