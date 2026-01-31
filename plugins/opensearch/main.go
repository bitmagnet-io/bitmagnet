package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/plugin-opensearch/client"
	"github.com/bitmagnet-io/plugin-opensearch/config"
	"github.com/bitmagnet-io/plugin-opensearch/i18n"
	"github.com/bitmagnet-io/plugin-opensearch/indexer"
	"github.com/bitmagnet-io/plugin-opensearch/search"
)

func main() {}

func init() {
	api.RegisterPlugin(&pluginAPI{})
}

type pluginAPI struct {
	config *config.Config
}

var _ api.Plugin = (*pluginAPI)(nil)

func (l *pluginAPI) Identify(ctx context.Context, localization *api.Empty) (*plugin.Identity, error) {
	return &plugin.Identity{
		Name:         "opensearch",
		Version:      "0.1.0",
		ConfigParams: config.Params,
	}, nil
}

func (l *pluginAPI) Localize(
	ctx context.Context,
	localizeParams *api.LocalizeParams,
) (*plugin.LocalizedContent, error) {
	localizer := i18n.NewLocalizer(localizeParams.GetAcceptLanguage())

	return &plugin.LocalizedContent{
		Description:  localizer.Localize("description"),
		ConfigParams: config.LocalizedContent(localizer),
	}, nil
}

func (l *pluginAPI) Configure(ctx context.Context, payload *api.JSONPayload) (*plugin.Contract, error) {
	var cfg config.Config
	err := json.Unmarshal(payload.Data, &cfg)
	if err != nil {
		return nil, err
	}

	l.config = &cfg

	return &plugin.Contract{
		Capabilities: &plugin.Capabilities{
			Indexer: &plugin.CapabilityIndexer{
				Name: "OpenSearch",
			},
			SearchAdapter: &plugin.CapabilitySearchAdapter{
				Name: "OpenSearch",
			},
		},
		Permissions: &plugin.Permissions{
			Http: &plugin.PermissionHTTP{
				Egress: []*http.Egress{
					{
						UrlPatterns: func() []string {
							patterns := make([]string, 0, len(l.config.Addresses))

							for _, addr := range l.config.Addresses {
								pattern := addr
								if !strings.HasSuffix(pattern, "/") {
									pattern += "/"
								}
								pattern += "**"

								patterns = append(patterns, pattern)
							}

							return patterns
						}(),
						Methods: []http.Method{
							http.Method_get,
							http.Method_post,
							http.Method_put,
							http.Method_delete,
							http.Method_head,
							http.Method_patch,
						},
					},
				},
			},
		},
	}, nil
}

func (l *pluginAPI) Instantiate(ctx context.Context, empty *api.Empty) (*api.Empty, error) {
	client, err := client.New(http_client.NewService(), *l.config)
	if err != nil {
		return nil, err
	}

	api.RegisterIndexer(indexer.New(client, *l.config))
	api.RegisterSearchAdapter(search.New(client, l.config.FullIndexPrefix()))

	return &api.Empty{}, nil
}
