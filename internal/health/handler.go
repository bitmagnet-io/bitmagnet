package health

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	HandlerConfig struct {
		statusCodeUp   int
		statusCodeDown int
		middleware     []Middleware
		resultWriter   ResultWriter
	}

	// Middleware is factory function that allows creating new instances of
	// a MiddlewareFunc. A MiddlewareFunc is expected to forward the function
	// call to the next MiddlewareFunc (passed in parameter 'next').
	// This way, a chain of interceptors is constructed that will eventually
	// invoke of the Checker.Check function. Each interceptor must therefore
	// invoke the 'next' interceptor. If the 'next' MiddlewareFunc is not called,
	// Checker.Check will never be executed.
	Middleware func(next MiddlewareFunc) MiddlewareFunc

	// MiddlewareFunc is a middleware for a health Handler (see NewHandler).
	// It is invoked each time an HTTP request is processed.
	MiddlewareFunc func(r *http.Request) CheckerResult

	// ResultWriter enabled a Handler (see NewHandler) to write the CheckerResult
	// to an http.ResponseWriter in a specific format. For example, the
	// JSONResultWriter writes the result in JSON format into the response body).
	ResultWriter interface {
		// Write writes a CheckerResult into a http.ResponseWriter in a format
		// that the ResultWriter supports (such as XML, JSON, etc.).
		// A ResultWriter is expected to write at least the following information into the http.ResponseWriter:
		// (1) A MIME type header (e.g., "Content-Type" : "application/json"),
		// (2) the HTTP status code that is passed in parameter statusCode (this is necessary due to ordering constraints
		// when writing into a http.ResponseWriter (see https://github.com/alexliesenfeld/health/issues/9), and
		// (3) the response body in the format that the ResultWriter supports.
		Write(result *CheckerResult, statusCode int, w http.ResponseWriter, r *http.Request) error
	}

	// JSONResultWriter writes a CheckerResult in JSON format into an
	// http.ResponseWriter. This ResultWriter is set by default.
	JSONResultWriter struct{}
)

// Write implements ResultWriter.Write.
func (rw *JSONResultWriter) Write(result *CheckerResult, statusCode int, w http.ResponseWriter, r *http.Request) error {
	jsonResp, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("cannot marshal response: %w", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResp)
	return err
}

// NewJSONResultWriter creates a new instance of a JSONResultWriter.
func NewJSONResultWriter() *JSONResultWriter {
	return &JSONResultWriter{}
}

// NewHandler creates a new health check http.Handler.
func NewHandler(checker Checker, options ...HandlerOption) http.HandlerFunc {
	cfg := createConfig(options)
	return func(w http.ResponseWriter, r *http.Request) {
		// Do the check (with configured middleware)
		result := withMiddleware(cfg.middleware, func(r *http.Request) CheckerResult {
			return checker.Check(r.Context())
		})(r)

		// Write HTTP response
		disableResponseCache(w)
		statusCode := mapHTTPStatusCode(result.Status, cfg.statusCodeUp, cfg.statusCodeDown)
		_ = cfg.resultWriter.Write(&result, statusCode, w, r)
	}
}

func disableResponseCache(w http.ResponseWriter) {
	// Avoid caching: https://www.ibm.com/garage/method/practices/manage/health-check-apis/
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
}

func mapHTTPStatusCode(status AvailabilityStatus, statusCodeUp int, statusCodeDown int) int {
	if status == StatusDown || status == StatusUnknown {
		return statusCodeDown
	}
	return statusCodeUp
}

func createConfig(options []HandlerOption) HandlerConfig {
	cfg := HandlerConfig{
		statusCodeDown: 503,
		statusCodeUp:   200,
		middleware:     []Middleware{},
	}

	for _, opt := range options {
		opt(&cfg)
	}

	if cfg.resultWriter == nil {
		cfg.resultWriter = &JSONResultWriter{}
	}

	return cfg
}

func withMiddleware(interceptors []Middleware, target MiddlewareFunc) MiddlewareFunc {
	chain := target
	for idx := len(interceptors) - 1; idx >= 0; idx-- {
		chain = interceptors[idx](chain)
	}
	return chain
}
