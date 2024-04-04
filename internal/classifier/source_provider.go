package classifier

import (
	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
	"os"
)

func newSourceProvider() sourceProvider {
	return mergeSourceProvider{
		providers: []sourceProvider{
			yamlSourceProvider{rawSourceProvider: coreSourceProvider{}},
			yamlSourceProvider{rawSourceProvider: xdgSourceProvider{}},
			yamlSourceProvider{rawSourceProvider: cwdSourceProvider{}},
		},
	}
}

type sourceProvider interface {
	source() (workflowSource, error)
}

type mergeSourceProvider struct {
	providers []sourceProvider
}

func (m mergeSourceProvider) source() (workflowSource, error) {
	source := workflowSource{}
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

func (y yamlSourceProvider) source() (workflowSource, error) {
	raw, err := y.rawSourceProvider.source()
	if err != nil {
		return workflowSource{}, err
	}
	rawWorkflow := make(map[string]interface{})
	parseErr := yaml.Unmarshal(raw, &rawWorkflow)
	if parseErr != nil {
		return workflowSource{}, parseErr
	}
	src := workflowSource{}
	decoder, decoderErr := newDecoder(&src)
	if decoderErr != nil {
		return workflowSource{}, decoderErr
	}
	if decodeErr := decoder.Decode(rawWorkflow); decodeErr != nil {
		return workflowSource{}, decodeErr
	}
	return src, nil
}

type coreSourceProvider struct{}

func (c coreSourceProvider) source() ([]byte, error) {
	return classifierCoreYaml, nil
}

type xdgSourceProvider struct{}

func (x xdgSourceProvider) source() ([]byte, error) {
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

func (x cwdSourceProvider) source() ([]byte, error) {
	if bytes, readErr := os.ReadFile("./classifier.yml"); readErr == nil {
		return bytes, nil
	} else if !os.IsNotExist(readErr) {
		return nil, readErr
	}
	return []byte{'{', '}'}, nil
}
