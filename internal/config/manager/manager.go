package manager

import (
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"gopkg.in/yaml.v3"
)

type FS fs.FS

type Manager struct {
	mtx        sync.Mutex
	fs         FS
	resolved   resolver.Resolved
	overridden ref.Map[yaml.Node]
	pending    ref.Set
}

func New(resolved resolver.Resolved, fs FS) *Manager {
	return &Manager{
		fs:         fs,
		resolved:   resolved,
		overridden: ref.NewMap[yaml.Node](),
		pending:    ref.NewSet(),
	}
}

func (m *Manager) Params() []*resolver.Param {
	return m.resolved.Params()
}

func (m *Manager) Validate(ref ref.Ref, value any) error {
	param, ok := m.resolved.Param(ref)
	if !ok {
		return fmt.Errorf("unknown parameter: %s", ref)
	}

	return param.ValidateAny(value)
}

func (m *Manager) Save(rf ref.Ref, value any) (*resolver.Param, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	param, ok := m.resolved.Param(rf)
	if !ok {
		return param, fmt.Errorf("unknown parameter: %s", rf)
	}

	decoded, err := param.DecodeYAMLAnyAny(value)
	if err != nil {
		return param, fmt.Errorf("failed to parse parameter value: %s: %w", rf, err)
	}

	err = param.ValidateAny(decoded)
	if err != nil {
		return nil, fmt.Errorf("value is invalid: %w", err)
	}

	yamlNode, err := param.EncodeYAMLAny(decoded)
	if err != nil {
		return param, fmt.Errorf("yaml encode failed: %w", err)
	}

	m.overridden.Set(rf, yamlNode)

	mapToWrite := ref.NewMap[yaml.Node]()
	for _, param := range m.resolved.Params() {
		if m.overridden.Has(param.Ref) || !shouldWriteSource(param.Source()) {
			continue
		}
		thisNode, err := param.EncodeYAMLAny(param.Value())
		if err != nil {
			return param, fmt.Errorf("yaml encode failed for key: %s: %w", param.Ref, err)
		}
		mapToWrite.Set(param.Ref, thisNode)
	}

	mapToWrite.SetAll(m.overridden)

	yamlBytes, err := yaml.Marshal(mapToWrite.StringMap())
	if err != nil {
		return param, fmt.Errorf("yaml marshal failed: %w", err)
	}

	err = m.fs.MkdirAll(".", 0o700)
	if err != nil {
		return param, fmt.Errorf("failed to create directory: %w", err)
	}

	err = m.fs.WriteFile(config.FilePersisted, yamlBytes, 0o600)
	if err != nil {
		return param, fmt.Errorf("failed to write file: %w", err)
	}

	if !param.IsDynamic() {
		m.pending.Set(rf, struct{}{})
	}

	return param, param.Save(decoded)
}

func (m *Manager) Delete(rf ref.Ref) (*resolver.Param, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	param, ok := m.resolved.Param(rf)
	if !ok {
		return nil, fmt.Errorf("unknown parameter: %s", rf)
	}

	mapToWrite := ref.NewMap[yaml.Node]()
	for _, param := range m.resolved.Params() {
		if m.overridden.Has(param.Ref) || !shouldWriteSource(param.Source()) || param.Ref.Equals(rf) {
			continue
		}
		thisNode, err := param.EncodeYAMLAny(param.Value())
		if err != nil {
			return nil, fmt.Errorf("yaml encode failed for key: %s: %w", param.Ref.String(), err)
		}
		mapToWrite.Set(param.Ref, thisNode)
	}

	m.overridden.Delete(rf)

	mapToWrite.SetAll(m.overridden)

	yamlBytes, err := yaml.Marshal(mapToWrite.StringMap())
	if err != nil {
		return nil, fmt.Errorf("yaml marshal failed: %w", err)
	}

	err = m.fs.MkdirAll(".", 0o700)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	err = m.fs.WriteFile(config.FilePersisted, yamlBytes, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	if !param.IsDynamic() {
		if param.Source() == config.SourcePersisted {
			m.pending.Set(rf, struct{}{})
		} else {
			m.pending.Delete(rf)
		}
	}

	if err = param.Delete(); err != nil {
		return nil, err
	}

	return param, nil
}

func (m *Manager) HasPending() bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.pending.Len() > 0
}

func shouldWriteSource(source string) bool {
	switch source {
	case config.SourceDynamic, config.SourcePending, config.SourcePersisted:
		return true
	default:
		return false
	}
}
