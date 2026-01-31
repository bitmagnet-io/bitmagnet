package transform_http

import (
	"net/http"

	proto "github.com/bitmagnet-io/bitmagnet/proto/common/http"
)

func HTTPMethodToProto(method string) proto.Method {
	switch method {
	case http.MethodGet:
		return proto.Method_get
	case http.MethodPost:
		return proto.Method_post
	case http.MethodPut:
		return proto.Method_put
	case http.MethodDelete:
		return proto.Method_delete
	case http.MethodPatch:
		return proto.Method_patch
	case http.MethodHead:
		return proto.Method_head
	case http.MethodOptions:
		return proto.Method_options
	default:
		return proto.Method_get
	}
}

func HTTPHeaderToProto(headers http.Header) map[string]string {
	result := make(map[string]string, len(headers))
	for k, v := range headers {
		result[k] = v[0]
	}

	return result
}

func HTTPMethodFromProto(method proto.Method) string {
	switch method {
	case proto.Method_get:
		return http.MethodGet
	case proto.Method_post:
		return http.MethodPost
	case proto.Method_put:
		return http.MethodPut
	case proto.Method_delete:
		return http.MethodDelete
	case proto.Method_patch:
		return http.MethodPatch
	case proto.Method_head:
		return http.MethodHead
	case proto.Method_options:
		return http.MethodOptions
	}

	return ""
}

func HTTPHeaderFromProto(headers map[string]string) http.Header {
	result := make(http.Header, len(headers))
	for k, v := range headers {
		result.Set(k, v)
	}

	return result
}
