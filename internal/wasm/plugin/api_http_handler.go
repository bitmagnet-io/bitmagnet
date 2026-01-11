package plugin

import (
	context "context"

	plugin_api "github.com/bitmagnet-io/bitmagnet/proto/api"
	http "github.com/bitmagnet-io/bitmagnet/proto/common/http"
)

type apiHTTPHandler struct {
	*apiService[plugin_api.HTTPHandler]
}

func (api *apiHTTPHandler) Config(ctx context.Context, empty *plugin_api.Empty) (*plugin_api.HTTPHandlerConfig, error) {
	var result *plugin_api.HTTPHandlerConfig

	err := api.do(ctx, func(httpHandler plugin_api.HTTPHandler) error {
		var err error
		result, err = httpHandler.Config(ctx, empty)
		return err
	})

	return result, err
}

func (api *apiHTTPHandler) HandleRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	var result *http.Response

	err := api.do(ctx, func(httpHandler plugin_api.HTTPHandler) error {
		var err error
		result, err = httpHandler.HandleRequest(ctx, request)
		return err
	})

	return result, err
}
