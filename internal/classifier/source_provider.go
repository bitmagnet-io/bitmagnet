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
			configSourceProvider{
				config:      config,
				tmdbEnabled: tmdbConfig.Enabled,
			},
		},
	}
}

type sourceProvider interface {
	source() (Source, error)
}

type mergeSourceProvider struct {
	providers []sourceProvider
}

func (m mergeSourceProvider) source() (Source, error) {
	source := Source{}
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

func (y yamlSourceProvider) source() (Source, error) {
	raw, err := y.rawSourceProvider.source()
	if err != nil {
		return Source{}, err
	}
	rawWorkflow := make(map[string]interface{})
	parseErr := yaml.Unmarshal(raw, &rawWorkflow)
	if parseErr != nil {
		return Source{}, parseErr
	}
	src := Source{}
	decoder, decoderErr := newDecoder(&src)
	if decoderErr != nil {
		return Source{}, decoderErr
	}
	if decodeErr := decoder.Decode(rawWorkflow); decodeErr != nil {
		return Source{}, decodeErr
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
	config      Config
	tmdbEnabled bool
}

func (c configSourceProvider) source() (Source, error) {
	fs := make(Flags)
	for k, v := range c.config.Flags {
		fs[k] = v
	}
	if c.config.DeleteXxx {
		fs["delete_xxx"] = true
	}
	if !c.tmdbEnabled {
		fs["tmdb_enabled"] = false
	}
	return Source{
		Keywords:   c.config.Keywords,
		Extensions: c.config.Extensions,
		Flags:      fs,
	}, nil
}
