package plugin

import (
	"context"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/search/adapter/multi"
	search_proto "github.com/bitmagnet-io/bitmagnet/internal/search/adapter/proto"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/target"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	proto_plugin "github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
)

type Plugin struct {
	ref          ref.Ref
	apiMutex     sync.Mutex
	api          api.Plugin
	newModule    func(context.Context) (*module, error)
	configParams []config_registry.Param
	httpEgress   *atomic.Value[[]*http.Egress]
}

var _ plugin.Plugin = (*Plugin)(nil)

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

func (p *Plugin) LocalizedContent(ctx context.Context, acceptLanguage ...string) (plugin.LocalizedContent, error) {
	lc, err := p.apiLocalize(ctx, acceptLanguage...)
	if err != nil {
		return plugin.LocalizedContent{}, err
	}

	configParams, err := slice.MapErr(
		lc.GetConfigParams(),
		func(l *proto_plugin.ConfigParamLocalizedContent) (plugin.LocalizedConfigParam, error) {
			ref, err := ref.Parse(p.ref.String() + "." + l.GetName())
			if err != nil {
				return plugin.LocalizedConfigParam{}, err
			}

			return plugin.LocalizedConfigParam{
				Ref:         ref,
				Description: l.GetDescription(),
			}, nil
		},
	)
	if err != nil {
		return plugin.LocalizedContent{}, err
	}

	return plugin.LocalizedContent{
		Ref:          p.ref,
		Description:  lc.GetDescription(),
		ConfigParams: configParams,
	}, nil
}

func (p *Plugin) apiLocalize(ctx context.Context, acceptLanguage ...string) (*proto_plugin.LocalizedContent, error) {
	p.apiMutex.Lock()
	defer p.apiMutex.Unlock()

	return p.api.Localize(ctx, &api.LocalizeParams{
		AcceptLanguage: acceptLanguage,
	})
}

func (p *Plugin) FXOption() fx.Option {
	instanceTag := `name:"` + p.ref.String() + `:instance"`

	options := []fx.Option{
		fx.Provide(
			fx.Annotate(
				func(env env.Env, cfg resolver.Resolved) (Instance, error) {
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

	options = append(options, fx.Provide(
		fx.Annotate(
			func(inst Instance) []indexer.Indexer {
				if contract := inst.Contract(); contract == nil ||
					contract.GetCapabilities() == nil ||
					contract.GetCapabilities().GetIndexer() == nil {
					return nil
				}

				return []indexer.Indexer{indexer.NewProto(inst.Indexer())}
			},
			fx.ParamTags(instanceTag),
			fx.ResultTags(`group:"indexers"`),
		),
	))

	options = append(options, fx.Provide(
		fx.Annotate(
			func(inst Instance) []multi.Index {
				contract := inst.Contract()
				if contract == nil ||
					contract.GetCapabilities() == nil ||
					contract.GetCapabilities().GetSearchAdapter() == nil {
					return nil
				}

				return []multi.Index{{
					Ref:     p.ref,
					Name:    contract.GetCapabilities().GetSearchAdapter().GetName(),
					Adapter: search_proto.New(inst.SearchAdapter()),
				}}
			},
			fx.ParamTags(instanceTag),
			fx.ResultTags(`group:"search_adapters"`),
		),
	))

	options = append(options, fx.Provide(
		fx.Annotate(
			func(inst Instance) []target.TorrentContentTarget {
				contract := inst.Contract()
				if contract == nil ||
					contract.GetCapabilities() == nil ||
					contract.GetCapabilities().GetTorrentTarget() == nil {
					return nil
				}

				return []target.TorrentContentTarget{
					target.NewTargetProto(
						p.ref,
						contract.GetCapabilities().GetTorrentTarget().GetName(),
						inst.TorrentTarget(),
					),
				}
			},
			fx.ParamTags(instanceTag),
			fx.ResultTags(`group:"torrent_targets"`),
		),
	))

	return fx.Options(options...)
}
