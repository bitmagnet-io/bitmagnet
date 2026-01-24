package batch_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler/testserver"
	batch "github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver/batch"
	"github.com/stretchr/testify/assert"
)

// revive:disable:line-length-limit
func TestNormalPOST(t *testing.T) {
	t.Parallel()

	h := testserver.New()
	h.AddTransport(batch.Transport{})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `{"query":"{ name }"}`)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"data":{"name":"test"}}`, resp.Body.String())
	})

	t.Run("decode failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", "notjson")
		assert.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"json body could not be decoded: invalid character 'o' in literal null (expecting 'u')"}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("parse failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `{"query": "!"}`)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"Unexpected !","locations":[{"line":1,"column":1}],"extensions":{"code":"GRAPHQL_PARSE_FAILED"}}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("validation failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `{"query": "{ title }"}`)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"Cannot query field \"title\" on type \"Query\".","locations":[{"line":1,"column":3}],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("invalid variable", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(
			h,
			"POST",
			"/graphql",
			`{"query": "query($id:Int!){find(id:$id)}","variables":{"id":false}}`,
		)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"cannot use bool as Int","path":["variable","id"],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("execution failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `{"query": "mutation { name }"}`)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"mutations are not supported"}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("validate content type", func(t *testing.T) {
		t.Parallel()

		doReq := func(
			handler http.Handler,
			method string,
			target string,
			body string,
			contentType string,
		) *httptest.ResponseRecorder {
			r := httptest.NewRequest(method, target, strings.NewReader(body))
			if contentType != "" {
				r.Header.Set("Content-Type", contentType)
			}

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			return w
		}

		validContentTypes := []string{
			"application/json",
			"application/json; charset=utf-8",
		}

		for _, contentType := range validContentTypes {
			t.Run(fmt.Sprintf("allow for content type %s", contentType), func(t *testing.T) {
				t.Parallel()

				resp := doReq(h, "POST", "/graphql", `{"query":"{ name }"}`, contentType)
				assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
				assert.JSONEq(t, `{"data":{"name":"test"}}`, resp.Body.String())
			})
		}

		invalidContentTypes := []string{
			"",
			"text/plain",

			// These content types are currently not supported, but they are supported by other GraphQL
			// servers, like express-graphql.
			"application/x-www-form-urlencoded",
			"application/graphql",
		}

		for _, tc := range invalidContentTypes {
			t.Run(fmt.Sprintf("reject for content type %s", tc), func(t *testing.T) {
				t.Parallel()

				resp := doReq(h, "POST", "/graphql", `{"query":"{ name }"}`, tc)
				assert.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
				assert.JSONEq(
					t,
					fmt.Sprintf(
						`{"errors":[{"message":"%s"}],"data":null}`,
						"transport not supported",
					),
					resp.Body.String(),
				)
			})
		}
	})
}

func TestBatchPOST(t *testing.T) {
	t.Parallel()

	h := testserver.New()
	h.AddTransport(batch.Transport{})

	t.Run("success1", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `[{"query":"{ name }"}]`)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `[{"data":{"name":"test"}}]`, resp.Body.String())
	})
	t.Run("success2", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `[{"query":"{ name }"},{"query":"{ name }"}]`)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `[{"data":{"name":"test"}},{"data":{"name":"test"}}]`, resp.Body.String())
	})
	t.Run("success3", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(
			h,
			"POST",
			"/graphql",
			`[{"query": "query($id:Int!){find(id:$id)}","variables":{"id":1}},{"query": "query($id:Int!){find(id:$id)}","variables":{"id":2}}]`,
		)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(t, `[{"data":{"name":"test"}},{"data":{"name":"test"}}]`, resp.Body.String())
	})
	t.Run("success and fail -> 200 OK", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(
			h,
			"POST",
			"/graphql",
			`[{"query": "query($id:Int!){find(id:$id)}","variables":{"id":1}},{"query": "query($id:Int!){find(id:$id)}","variables":{"id":false}}]`,
		)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`[{"data":{"name":"test"}},{"errors":[{"message":"cannot use bool as Int","path":["variable","id"],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}}],"data":null}]`,
			resp.Body.String(),
		)
	})

	t.Run("must not be an empty array", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", "[]")
		assert.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"json body must not be an empty array"}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("decode failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", "notjson")
		assert.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`{"errors":[{"message":"json body could not be decoded: invalid character 'o' in literal null (expecting 'u')"}],"data":null}`,
			resp.Body.String(),
		)
	})

	t.Run("parse failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `[{"query": "!"}]`)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`[{"errors":[{"message":"Unexpected !","locations":[{"line":1,"column":1}],"extensions":{"code":"GRAPHQL_PARSE_FAILED"}}],"data":null}]`,
			resp.Body.String(),
		)
	})

	t.Run("validation failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `[{"query": "{ title }"}]`)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`[{"errors":[{"message":"Cannot query field \"title\" on type \"Query\".","locations":[{"line":1,"column":3}],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}}],"data":null}]`,
			resp.Body.String(),
		)
	})

	t.Run("invalid variable", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(
			h,
			"POST",
			"/graphql",
			`[{"query": "query($id:Int!){find(id:$id)}","variables":{"id":false}}]`,
		)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`[{"errors":[{"message":"cannot use bool as Int","path":["variable","id"],"extensions":{"code":"GRAPHQL_VALIDATION_FAILED"}}],"data":null}]`,
			resp.Body.String(),
		)
	})

	t.Run("execution failure", func(t *testing.T) {
		t.Parallel()

		resp := doRequest(h, "POST", "/graphql", `[{"query": "mutation { name }"}]`)
		assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
		assert.JSONEq(
			t,
			`[{"errors":[{"message":"mutations are not supported"}],"data":null}]`,
			resp.Body.String(),
		)
	})

	t.Run("validate content type", func(t *testing.T) {
		t.Parallel()

		doReq := func(handler http.Handler, method string, target string, body string, contentType string) *httptest.ResponseRecorder {
			r := httptest.NewRequest(method, target, strings.NewReader(body))
			if contentType != "" {
				r.Header.Set("Content-Type", contentType)
			}

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			return w
		}

		validContentTypes := []string{
			"application/json",
			"application/json; charset=utf-8",
		}

		for _, contentType := range validContentTypes {
			t.Run(fmt.Sprintf("allow for content type %s", contentType), func(t *testing.T) {
				t.Parallel()

				resp := doReq(h, "POST", "/graphql", `[{"query":"{ name }"}]`, contentType)
				assert.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
				assert.JSONEq(t, `[{"data":{"name":"test"}}]`, resp.Body.String())
			})
		}

		invalidContentTypes := []string{
			"",
			"text/plain",

			// These content types are currently not supported, but they are supported by other GraphQL
			// servers, like express-graphql.
			"application/x-www-form-urlencoded",
			"application/graphql",
		}

		for _, tc := range invalidContentTypes {
			t.Run(fmt.Sprintf("reject for content type %s", tc), func(t *testing.T) {
				t.Parallel()

				resp := doReq(h, "POST", "/graphql", `{"query":"{ name }"}`, tc)
				assert.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
				assert.JSONEq(
					t,
					fmt.Sprintf(
						`{"errors":[{"message":"%s"}],"data":null}`,
						"transport not supported",
					),
					resp.Body.String(),
				)
			})
		}
	})
}

func doRequest(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	return w
}
