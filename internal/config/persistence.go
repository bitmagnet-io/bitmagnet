package config

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/config/lookup"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"gopkg.in/yaml.v3"
)

func (m *Manager) Save(rf ref.Ref, value any) (Param, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	param, ok := m.resolved.Param(rf)
	if !ok {
		return Param{}, fmt.Errorf("unknown parameter: %s", rf)
	}

	decoded, err := param.DecodeYAMLAnyAny(value)
	if err != nil {
		return Param{}, fmt.Errorf("failed to parse parameter value: %s: %w", rf, err)
	}

	err = param.ValidateAny(decoded)
	if err != nil {
		return Param{}, fmt.Errorf("value is invalid: %w", err)
	}

	yamlNode, err := param.EncodeYAMLAny(decoded)
	if err != nil {
		return Param{}, fmt.Errorf("yaml encode failed: %w", err)
	}

	m.overridden.Set(rf, yamlNode)

	mapToWrite := ref.NewMap[yaml.Node]()

	for _, param := range m.resolved.Params() {
		if m.overridden.Has(param.Ref) || !shouldWriteSource(param.Source()) {
			continue
		}

		thisNode, err := param.EncodeYAMLAny(param.Value())
		if err != nil {
			return Param{}, fmt.Errorf("yaml encode failed for key: %s: %w", param.Ref, err)
		}

		mapToWrite.Set(param.Ref, thisNode)
	}

	mapToWrite.SetAll(m.overridden)

	yamlBytes, err := yaml.Marshal(mapToWrite.StringMap())
	if err != nil {
		return Param{}, fmt.Errorf("yaml marshal failed: %w", err)
	}

	err = m.fs.MkdirAll(".", 0o700)
	if err != nil {
		return Param{}, fmt.Errorf("failed to create directory: %w", err)
	}

	err = m.fs.WriteFile(lookup.FilePersisted, yamlBytes, 0o600)
	if err != nil {
		return Param{}, fmt.Errorf("failed to write file: %w", err)
	}

	if !param.IsDynamic() {
		m.pending.Set(rf, struct{}{})
	}

	return transformParam(param), param.Save(decoded)
}

func (m *Manager) Delete(rf ref.Ref) (Param, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	param, ok := m.resolved.Param(rf)
	if !ok {
		return Param{}, fmt.Errorf("unknown parameter: %s", rf)
	}

	mapToWrite := ref.NewMap[yaml.Node]()

	for _, param := range m.resolved.Params() {
		if m.overridden.Has(param.Ref) || !shouldWriteSource(param.Source()) || param.Ref.Equals(rf) {
			continue
		}

		thisNode, err := param.EncodeYAMLAny(param.Value())
		if err != nil {
			return Param{}, fmt.Errorf("yaml encode failed for key: %s: %w", param.Ref, err)
		}

		mapToWrite.Set(param.Ref, thisNode)
	}

	m.overridden.Delete(rf)

	mapToWrite.SetAll(m.overridden)

	yamlBytes, err := yaml.Marshal(mapToWrite.StringMap())
	if err != nil {
		return Param{}, fmt.Errorf("yaml marshal failed: %w", err)
	}

	err = m.fs.MkdirAll(".", 0o700)
	if err != nil {
		return Param{}, fmt.Errorf("failed to create directory: %w", err)
	}

	err = m.fs.WriteFile(lookup.FilePersisted, yamlBytes, 0o600)
	if err != nil {
		return Param{}, fmt.Errorf("failed to write file: %w", err)
	}

	if !param.IsDynamic() {
		if param.Source() == lookup.SourcePersisted {
			m.pending.Set(rf, struct{}{})
		} else {
			m.pending.Delete(rf)
		}
	}

	if err = param.Delete(); err != nil {
		return Param{}, err
	}

	return transformParam(param), nil
}

func shouldWriteSource(source string) bool {
	switch source {
	case lookup.SourceDynamic, lookup.SourcePending, lookup.SourcePersisted:
		return true
	default:
		return false
	}
}
