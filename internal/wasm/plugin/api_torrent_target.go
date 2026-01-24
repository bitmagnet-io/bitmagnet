package plugin

import (
	"context"

	plugin_api "github.com/bitmagnet-io/bitmagnet/proto/api"
)

type apiTorrentTarget struct {
	*apiService[plugin_api.TorrentTarget]
}

func (api *apiTorrentTarget) DataSchema(ctx context.Context, _ *plugin_api.Empty) (*plugin_api.JSONPayload, error) {
	var result *plugin_api.JSONPayload

	err := api.do(ctx, func(torrentTarget plugin_api.TorrentTarget) error {
		var err error

		result, err = torrentTarget.DataSchema(ctx, &plugin_api.Empty{})

		return err
	})

	return result, err
}

func (api *apiTorrentTarget) UISchema(
	ctx context.Context,
	params *plugin_api.SendTorrentsUISchemaParams,
) (*plugin_api.JSONPayload, error) {
	var result *plugin_api.JSONPayload

	err := api.do(ctx, func(torrentTarget plugin_api.TorrentTarget) error {
		var err error

		result, err = torrentTarget.UISchema(ctx, params)

		return err
	})

	return result, err
}

func (api *apiTorrentTarget) Send(
	ctx context.Context,
	params *plugin_api.SendTorrentsParams,
) (*plugin_api.JSONPayload, error) {
	var result *plugin_api.JSONPayload

	err := api.do(ctx, func(torrentTarget plugin_api.TorrentTarget) error {
		var err error

		result, err = torrentTarget.Send(ctx, params)

		return err
	})

	return result, err
}
