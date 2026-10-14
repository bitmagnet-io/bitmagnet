package tmdb

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequesterLogger_RedactsAPIKeyOnError(t *testing.T) {
	t.Parallel()

	const apiKey = "super-secret-api-key"

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core).Sugar()

	r := requesterLogger{
		requester: requester{
			resty: resty.New().
				SetBaseURL(srv.URL).
				SetQueryParam("api_key", apiKey).
				SetTimeout(50 * time.Millisecond),
		},
		logger: logger,
	}

	_, err := r.Request(context.Background(), "/authentication", nil, nil)
	if err == nil {
		t.Fatal("expected request to fail")
	}

	entries := logs.All()
	if len(entries) == 0 {
		t.Fatal("expected an error log entry")
	}

	for _, e := range entries {
		if strings.Contains(e.Message, apiKey) {
			t.Fatalf("api key leaked in log message: %q", e.Message)
		}

		for _, f := range e.Context {
			serialized := f.Key + f.String + fmt.Sprintf("%v", f.Interface)
			if strings.Contains(serialized, apiKey) {
				t.Fatalf("api key leaked in log field %q", f.Key)
			}
		}
	}
}
