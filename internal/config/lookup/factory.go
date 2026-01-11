package lookup

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"gopkg.in/yaml.v3"
)

func NewFromEnv(env env.Env) (Lookup, error) {
	var lookups LookupChain

	persistedPath := filepath.Join(config.SubpathPersisted, config.FilePersisted)
	if yamlBytes, err := env.FSData().ReadFile(persistedPath); err == nil {
		lookup, err := NewFromYAMLFlat(config.SourcePersisted, yamlBytes)
		if err != nil {
			return nil, err
		}
		lookups = append(lookups, lookup)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %w: %s: %w", Err, ErrReadFile, persistedPath, err)
	}

	lookups = append(lookups, LookupEnv("env", env))

	if yamlBytes, err := env.FSCurrent().ReadFile("./config.yml"); err == nil {
		lookup, err := NewFromYAML("file:./config.yml", yamlBytes)
		if err != nil {
			return nil, err
		}
		lookups = append(lookups, lookup)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %w: %s: %w", Err, ErrReadFile, "./config.yml", err)
	}

	xdgFile := filepath.Join("bitmagnet", "config.yml")
	if yamlBytes, err := env.FSConfig().ReadFile(xdgFile); err == nil {
		lookup, err := NewFromYAML("file:"+xdgFile, yamlBytes)
		if err != nil {
			return nil, err
		}
		lookups = append(lookups, lookup)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %w: %s: %w", Err, ErrReadFile, xdgFile, err)
	}

	return lookups, nil
}

func NewFromYAML(source string, yamlBytes []byte) (Lookup, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &node); err != nil {
		return nil, fmt.Errorf("%w: %w: %s: %w", Err, ErrYAMLUnmarshal, source, err)
	}

	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = *node.Content[0]
	}

	return lookupYAML{source: source, node: node}, nil
}

func NewFromYAMLFlat(source string, yamlBytes []byte) (Lookup, error) {
	var nodeMap map[string]yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &nodeMap); err != nil {
		return nil, fmt.Errorf("%w: %w: %s: %w", Err, ErrYAMLUnmarshal, source, err)
	}

	return lookupYAMLFlat{source: source, nodeMap: nodeMap}, nil
}
