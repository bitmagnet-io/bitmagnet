package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	"github.com/bitmagnet-io/plugin-qbittorrent/config"
	"github.com/bitmagnet-io/plugin-qbittorrent/i18n"
	"github.com/bitmagnet-io/plugin-qbittorrent/target"
)

func main() {}

func init() {
	api.RegisterPlugin(&pluginAPI{})
}

type pluginAPI struct {
	config *config.Config
}

var _ api.Plugin = (*pluginAPI)(nil)

func (l *pluginAPI) Identify(ctx context.Context, _ *api.Empty) (*plugin.Identity, error) {
	return &plugin.Identity{
		Name:         "qbittorrent",
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
			TorrentTarget: &plugin.CapabilityTorrentTarget{
				Name: "qBittorrent",
			},
		},
		Permissions: &plugin.Permissions{
			Http: &plugin.PermissionHTTP{
				Egress: []*http.Egress{
					{
						UrlPatterns: func() []string {
							pattern := cfg.URL
							if !strings.HasSuffix(pattern, "/") {
								pattern += "/"
							}
							pattern += "**"

							return []string{pattern}
						}(),
						Methods: []http.Method{http.Method_get, http.Method_post},
					},
				},
			},
		},
	}, nil
}

func (l *pluginAPI) Instantiate(ctx context.Context, empty *api.Empty) (*api.Empty, error) {
	api.RegisterTorrentTarget(target.New(*l.config))

	return &api.Empty{}, nil
}
