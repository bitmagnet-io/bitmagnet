package classifier

import (
	"github.com/adrg/xdg"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"gopkg.in/yaml.v3"
	"os"
)

func newSourceProvider(config Config, tmdbConfig tmdb.Config) sourceProvider {
	return mergeSourceProvider{
		providers: []sourceProvider{
			yamlSourceProvider{rawSourceProvider: coreSourceProvider{}},
			yamlSourceProvider{rawSourceProvider: xdgSourceProvider{}},
			yamlSourceProvider{rawSourceProvider: cwdSourceProvider{}},
			configSourceProvider{config: config},
			disableTmdbSourceProvider{enabled: tmdbConfig.Enabled},
		},
	}
}

type sourceProvider interface {
	source() (WorkflowSource, error)
}

type mergeSourceProvider struct {
	providers []sourceProvider
}

func (m mergeSourceProvider) source() (WorkflowSource, error) {
	source := WorkflowSource{}
	for _, p := range m.providers {
		s, err := p.source()
		if err != nil {
			return source, err
		}
		if merged, err := source.merge(s); err != nil {
			return source, err
		} else {
			source = merged
		}
	}
	return source, nil
}

type rawSourceProvider interface {
	source() ([]byte, error)
}

type yamlSourceProvider struct {
	rawSourceProvider
}

func (y yamlSourceProvider) source() (WorkflowSource, error) {
	raw, err := y.rawSourceProvider.source()
	if err != nil {
		return WorkflowSource{}, err
	}
	rawWorkflow := make(map[string]interface{})
	parseErr := yaml.Unmarshal(raw, &rawWorkflow)
	if parseErr != nil {
		return WorkflowSource{}, parseErr
	}
	src := WorkflowSource{}
	decoder, decoderErr := newDecoder(&src)
	if decoderErr != nil {
		return WorkflowSource{}, decoderErr
	}
	if decodeErr := decoder.Decode(rawWorkflow); decodeErr != nil {
		return WorkflowSource{}, decodeErr
	}
	return src, nil
}

type coreSourceProvider struct{}

func (c coreSourceProvider) source() ([]byte, error) {
	return classifierCoreYaml, nil
}

type xdgSourceProvider struct{}

func (_ xdgSourceProvider) source() ([]byte, error) {
	if path, pathErr := xdg.ConfigFile("bitmagnet/classifier.yml"); pathErr == nil {
		if bytes, readErr := os.ReadFile(path); readErr == nil {
			return bytes, nil
		} else if !os.IsNotExist(readErr) {
			return nil, readErr
		}
	}
	return []byte{'{', '}'}, nil
}

type cwdSourceProvider struct{}

func (_ cwdSourceProvider) source() ([]byte, error) {
	if bytes, readErr := os.ReadFile("./classifier.yml"); readErr == nil {
		return bytes, nil
	} else if !os.IsNotExist(readErr) {
		return nil, readErr
	}
	return []byte{'{', '}'}, nil
}

type configSourceProvider struct {
	config Config
}

func (c configSourceProvider) source() (WorkflowSource, error) {
	return WorkflowSource{
		Keywords:   c.config.Keywords,
		Extensions: c.config.Extensions,
		Flags:      c.config.Flags,
	}, nil
}

type disableTmdbSourceProvider struct {
	enabled bool
}

func (d disableTmdbSourceProvider) source() (WorkflowSource, error) {
	if !d.enabled {
		return WorkflowSource{
			Flags: map[string]any{
				"tmdb_enabled": false,
			},
		}, nil
	}
	return WorkflowSource{}, nil
}
