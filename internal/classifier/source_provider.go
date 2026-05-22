package classifier

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"gopkg.in/yaml.v3"
)

func newSourceProvider(config Config, tmdbConfig tmdb.Config) sourceProvider {
	var providers []sourceProvider
	providers = append(providers, coreSourceProvider{}.provider())
	providers = append(providers, xdgSourceProvider{}.providers()...)
	providers = append(providers, cwdSourceProvider{}.providers()...)
	providers = append(providers, extraSourceProvider{}.providers()...)
	providers = append(providers, configSourceProvider{
		config:      config,
		tmdbEnabled: tmdbConfig.Enabled,
	})

	return mergeSourceProvider{providers: providers}
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

		merged, err := source.merge(s)
		if err != nil {
			return source, err
		}

		source = merged
	}

	return source, nil
}

type yamlSourceProvider struct {
	raw []byte
	err error
}

func (y yamlSourceProvider) source() (Source, error) {
	if y.err != nil {
		return Source{}, y.err
	}

	rawWorkflow := make(map[string]interface{})

	parseErr := yaml.Unmarshal(y.raw, &rawWorkflow)
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

func (coreSourceProvider) provider() sourceProvider {
	return yamlSourceProvider{raw: classifierCoreYaml}
}

type xdgSourceProvider struct{}

func (yamlSourceProvider) providers(path string, wildcard bool) []sourceProvider {
	glob := path

	if wildcard {
		dir, fname := filepath.Split(path)
		glob = dir + "classifier*" + filepath.Ext(fname)
	}

	paths, err := filepath.Glob(glob)
	if err != nil {
		return []sourceProvider{yamlSourceProvider{err: err}}
	}

	providers := make([]sourceProvider, len(paths))

	for i, path := range paths {
		bytes, readErr := os.ReadFile(path)
		providers[i] = yamlSourceProvider{raw: bytes, err: readErr}
	}

	return providers
}

func (xdgSourceProvider) providers() []sourceProvider {
	path, err := xdg.ConfigFile("bitmagnet/classifier.yml")
	if err != nil {
		return []sourceProvider{yamlSourceProvider{err: err}}
	}

	return yamlSourceProvider{}.providers(path, true)
}

type extraSourceProvider struct{}

func (extraSourceProvider) providers() []sourceProvider {
	var extraConfigFiles []string

	for _, rawEnvEntry := range os.Environ() {
		parts := strings.SplitN(rawEnvEntry, "=", 2)
		if parts[0] == extraFilesKey {
			extraConfigFiles = strings.Split(parts[1], ",")
		}
	}

	var providers []sourceProvider
	for _, path := range extraConfigFiles {
		providers = append(providers, yamlSourceProvider{}.providers(path, false)...)
	}

	return providers
}

type cwdSourceProvider struct{}

func (cwdSourceProvider) providers() []sourceProvider {
	return yamlSourceProvider{}.providers("./classifier.yml", true)
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

const extraFilesKey = "EXTRA_CLASSIFIER_FILES"
