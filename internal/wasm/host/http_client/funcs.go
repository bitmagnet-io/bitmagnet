package http_client

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/bitmagnet/proto/transformer/transform_http"
	"github.com/go-resty/resty/v2"
	"github.com/tetratelabs/wazero"
)

func Instantiator() func(ctx context.Context, runtime wazero.Runtime) error {
	client := client{
		resty: resty.New(),
	}

	return func(ctx context.Context, runtime wazero.Runtime) error {
		return http_client.Instantiate(ctx, runtime, client)
	}
}

type client struct {
	resty *resty.Client
}

var _ http_client.Service = client{}

func (c client) Request(ctx context.Context, req *http.Request) (*http.Response, error) {
	res, err := c.resty.R().
		SetContext(ctx).
		SetHeaders(req.GetHeaders()).
		SetBody(req.GetBody()).
		Execute(transform_http.HTTPMethodFromProto(req.GetMethod()), req.GetUrl())
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Status:  int32(res.StatusCode()),
		Headers: transform_http.HTTPHeaderToProto(res.Header()),
		Body:    res.Body(),
	}, nil
}
