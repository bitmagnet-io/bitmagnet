package httpclient

import "net/http"

type TransportDecorator func(t http.RoundTripper) http.RoundTripper
