package manager

import (
	"fmt"
	"maps"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/fs"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"gopkg.in/yaml.v3"
)

type FS fs.FS

type Manager struct {
	mtx        sync.Mutex
	fs         FS
	resolved   resolver.Resolved
	overridden map[string]yaml.Node
	pending    map[string]struct{}
}

func New(resolved resolver.Resolved, fs FS) *Manager {
	return &Manager{
		fs:         fs,
		resolved:   resolved,
		overridden: make(map[string]yaml.Node),
		pending:    make(map[string]struct{}),
	}
}

func (m *Manager) Params() []*resolver.Param {
	return m.resolved.Params()
}

func (m *Manager) Save(ref ref.Ref, value string) (*resolver.Param, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	param, ok := m.resolved.Param(ref)
	if !ok {
		return param, fmt.Errorf("unknown parameter: %s", ref)
	}

	parsed, err := param.ParseAny(value)
	if err != nil {
		return param, fmt.Errorf("failed to parse parameter value: %s: %w", ref, err)
	}

	yamlNode, err := param.EncodeYAMLAny(parsed)
	if err != nil {
		return param, fmt.Errorf("yaml encode failed: %w", err)
	}

	m.overridden[ref.String()] = yamlNode

	mapToWrite := make(map[string]yaml.Node)
	for _, param := range m.resolved.Params() {
		_, isOverridden := m.overridden[ref.String()]
		if isOverridden || !shouldWriteSource(param.Source()) {
			continue
		}
		thisNode, err := param.EncodeYAMLAny(param.Value())
		if err != nil {
			return param, fmt.Errorf("yaml encode failed for key: %s: %w", param.Ref.String(), err)
		}
		mapToWrite[param.Ref.String()] = thisNode
	}

	maps.Copy(mapToWrite, m.overridden)

	yamlBytes, err := yaml.Marshal(mapToWrite)
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
		m.pending[ref.String()] = struct{}{}
	}

	param.Save(parsed)

	return param, nil
}

func (m *Manager) Delete(ref ref.Ref) (*resolver.Param, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	param, ok := m.resolved.Param(ref)
	if !ok {
		return nil, fmt.Errorf("unknown parameter: %s", ref)
	}

	mapToWrite := make(map[string]yaml.Node)
	for _, param := range m.resolved.Params() {
		_, isOverridden := m.overridden[ref.String()]
		if isOverridden || !shouldWriteSource(param.Source()) || param.Ref.String() == ref.String() {
			continue
		}
		thisNode, err := param.EncodeYAMLAny(param.Value())
		if err != nil {
			return nil, fmt.Errorf("yaml encode failed for key: %s: %w", param.Ref.String(), err)
		}
		mapToWrite[param.Ref.String()] = thisNode
	}

	delete(m.overridden, ref.String())

	maps.Copy(mapToWrite, m.overridden)

	yamlBytes, err := yaml.Marshal(mapToWrite)
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

	param.Delete()

	if !param.IsDynamic() {
		if param.Source() == config.SourcePersisted {
			m.pending[ref.String()] = struct{}{}
		} else {
			delete(m.pending, ref.String())
		}
	}

	return param, nil
}

func (m *Manager) HasPending() bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return len(m.pending) > 0
}

func shouldWriteSource(source string) bool {
	switch source {
	case config.SourceDynamic, config.SourcePending, config.SourcePersisted:
		return true
	default:
		return false
	}
}
