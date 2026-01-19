package plugin

import (
	"fmt"

	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/search/adapter/multi"
	"github.com/bitmagnet-io/bitmagnet/internal/search/adapter/proto"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tetratelabs/wazero"
	"go.uber.org/fx"
)

type Plugin struct {
	ref              ref.Ref
	manifest         Manifest
	configParams     []config_registry.Param
	compilationCache wazero.CompilationCache
	data             []byte
}

func (p *Plugin) Ref() ref.Ref {
	return p.ref
}

func (*Plugin) ActivationRef() ref.Nullable {
	return ref.Nullable{}
}

func (*Plugin) Dependencies() []ref.Ref {
	return nil
}

func (p *Plugin) ConfigParams() []config_registry.Param {
	return p.configParams
}

func (*Plugin) Errors() ref.Map[error] {
	return ref.NewMap[error]()
}

func (*Plugin) I18nMessages() []*i18n.Message {
	return nil
}

func (*Plugin) Commands() []plugin.Command {
	return nil
}

func (p *Plugin) FXOption() fx.Option {
	type instance Instance

	instanceTag := fmt.Sprintf(`name:"%v"`, p.Ref())
	options := []fx.Option{
		fx.Provide(
			fx.Annotate(
				func(env env.Env, cfg resolver.Resolved) (instance, error) {
					instance, err := p.NewInstance(env, cfg)
					if err != nil {
						return nil, err
					}

					return instance, nil
				},
				fx.ResultTags(instanceTag),
			),
		),
	}

	if p.manifest.Capabilities.Indexer != nil {
		options = append(options, fx.Provide(
			fx.Annotate(
				func(inst instance) indexer.Indexer {
					return indexer.NewProto(inst.Indexer())
				},
				fx.ResultTags(`group:"indexers"`),
				fx.ParamTags(instanceTag),
			),
		))
	}

	if caplty := p.manifest.Capabilities.SearchAdapter; caplty != nil {
		options = append(options, fx.Provide(
			fx.Annotate(
				func(inst instance) multi.Index {
					return multi.Index{
						Ref:     p.ref,
						Name:    caplty.Name,
						Adapter: proto.New(inst.SearchAdapter()),
					}
				},
				fx.ResultTags(`group:"search_adapters"`),
				fx.ParamTags(instanceTag),
			),
		))
	}

	return fx.Options(options...)
}
