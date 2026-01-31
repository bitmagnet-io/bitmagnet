package config

import (
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"gopkg.in/yaml.v3"
)

type FS fs.FS

type Manager struct {
	mtx        sync.Mutex
	fs         FS
	resolved   resolver.Resolved
	overridden ref.Map[yaml.Node]
	pending    ref.Set
	i18n       *i18n.Bundle
}

func New(resolved resolver.Resolved, fs FS, i18n *i18n.Bundle) *Manager {
	return &Manager{
		fs:         fs,
		resolved:   resolved,
		overridden: ref.NewMap[yaml.Node](),
		pending:    ref.NewSet(),
		i18n:       i18n,
	}
}

func (m *Manager) Params() []Param {
	return slice.Map(m.resolved.Params(), transformParam)
}

func (m *Manager) Validate(ref ref.Ref, value any) error {
	param, ok := m.resolved.Param(ref)
	if !ok {
		return fmt.Errorf("unknown parameter: %s", ref)
	}

	return param.ValidateAny(value)
}

func (m *Manager) HasPending() bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.pending.Len() > 0
}
