package client

import (
	"bytes"
	"fmt"
	"io"
	nethttp "net/http"

	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/bitmagnet/proto/transformer/transform_http"
)

type transport struct {
	httpClient http_client.Service
}

func (t *transport) RoundTrip(req *nethttp.Request) (*nethttp.Response, error) {
	ctx := req.Context()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	_ = req.Body.Close()

	res, err := t.httpClient.Request(ctx, &http.Request{
		Method:  transform_http.HTTPMethodToProto(req.Method),
		Url:     req.URL.String(),
		Headers: transform_http.HTTPHeaderToProto(req.Header),
		Body:    body,
	})
	if err != nil {
		return nil, err
	}

	return &nethttp.Response{
		StatusCode: int(res.GetStatus()),
		Header:     transform_http.HTTPHeaderFromProto(res.GetHeaders()),
		Body:       io.NopCloser(bytes.NewReader(res.GetBody())),
		Request:    req,
	}, nil
}
