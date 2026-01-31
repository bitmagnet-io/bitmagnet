package http_client

import (
	"context"
	"fmt"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/bitmagnet/proto/transformer/transform_http"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-resty/resty/v2"
	"github.com/tetratelabs/wazero"
)

func Instantiator(egress *atomic.Value[[]*http.Egress]) func(ctx context.Context, runtime wazero.Runtime) error {
	client := client{
		resty: resty.New(),
		checker: atomic.MapReader(egress, func(egress []*http.Egress) checkerFunc {
			return func(req *http.Request) bool {
				for _, e := range egress {
					if !slices.Contains(e.GetMethods(), req.GetMethod()) {
						continue
					}

					for _, pattern := range e.GetUrlPatterns() {
						if match, _ := doublestar.PathMatch(
							pattern,
							req.GetUrl(),
						); match {
							return true
						}
					}
				}

				return false
			}
		}),
	}

	return func(ctx context.Context, runtime wazero.Runtime) error {
		return http_client.Instantiate(ctx, runtime, client)
	}
}

type client struct {
	resty   *resty.Client
	checker atomic.Reader[checkerFunc]
}

var _ http_client.Service = client{}

func (c client) Request(ctx context.Context, req *http.Request) (*http.Response, error) {
	method := transform_http.HTTPMethodFromProto(req.GetMethod())

	checker := c.checker.Get()
	if !checker(req) {
		return nil, fmt.Errorf("request not allowed: %s %s", method, req.GetUrl())
	}

	res, err := c.resty.R().
		SetContext(ctx).
		SetHeaders(req.GetHeaders()).
		SetBody(req.GetBody()).
		Execute(method, req.GetUrl())
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Status:  int32(res.StatusCode()),
		Headers: transform_http.HTTPHeaderToProto(res.Header()),
		Body:    res.Body(),
	}, nil
}

type checkerFunc func(*http.Request) bool
