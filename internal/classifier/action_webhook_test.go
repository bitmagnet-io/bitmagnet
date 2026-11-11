package classifier

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/stretchr/testify/assert"
)

type rawTestSourceProvider struct {
	raw []byte
}

func (p rawTestSourceProvider) source() ([]byte, error) {
	return p.raw, nil
}

func TestWebhookAction(t *testing.T) {
	t.Parallel()

	t.Run("webhook classifies a torrent", func(t *testing.T) {
		t.Parallel()

		var received webhookRequest

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"matched":     true,
				"contentType": "movie",
				"baseTitle":   "The Webhook Movie",
				"year":        2000,
				"languages":   []string{"en"},
				"content": map[string]any{
					"source": "webhook",
					"id":     "456",
					"title":  "The Webhook Movie",
					"type":   "movie",
				},
			})
		}))
		defer server.Close()

		source, err := yamlSourceProvider{rawSourceProvider: rawTestSourceProvider{raw: []byte(`
workflows:
  default:
    - webhook:
        url: "` + server.URL + `"
        headers:
          X-Test: "1"
`)}}.source()
		if err != nil {
			t.Fatal(err)
		}

		mocks := newTestClassifierMocks(t)
		workflow, err := mocks.compiler.Compile(source)
		if err != nil {
			t.Fatal(err)
		}

		result, err := workflow.Run(context.Background(), "default", Flags{}, model.Torrent{
			Name:        "The Webhook Movie (2000).mkv",
			FilesStatus: model.FilesStatusSingle,
			Extension:   model.NewNullString("mkv"),
			Size:        1000000000,
		})
		if err != nil {
			t.Fatal(err)
		}

		// verify the request the webhook received
		assert.Equal(t, "The Webhook Movie (2000).mkv", received.Torrent.Name)
		assert.Equal(t, "The Webhook Movie (2000)", received.Torrent.BaseName)
		assert.Equal(t, "mkv", received.Torrent.Files[0].Extension)

		// verify the result the webhook produced
		assert.Equal(t, model.NewNullContentType(model.ContentTypeMovie), result.ContentType)
		assert.Equal(t, model.NewNullString("The Webhook Movie"), result.BaseTitle)
		assert.Equal(t, model.Date{Year: 2000}, result.Date)
		_, hasEn := result.Languages[model.Language("en")]
		assert.True(t, hasEn)
		assert.NotNil(t, result.Content)
		assert.Equal(t, "webhook", result.Content.Source)
		assert.Equal(t, "456", result.Content.ID)
	})

	t.Run("webhook unmatched", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"matched": false,
			})
		}))
		defer server.Close()

		source, err := yamlSourceProvider{rawSourceProvider: rawTestSourceProvider{raw: []byte(`
workflows:
  default:
    - webhook:
        url: "` + server.URL + `"
`)}}.source()
		if err != nil {
			t.Fatal(err)
		}

		mocks := newTestClassifierMocks(t)
		workflow, err := mocks.compiler.Compile(source)
		if err != nil {
			t.Fatal(err)
		}

		_, err = workflow.Run(context.Background(), "default", Flags{}, model.Torrent{
			Name:        "Unknown Thing.mkv",
			FilesStatus: model.FilesStatusSingle,
			Extension:   model.NewNullString("mkv"),
			Size:        1000000000,
		})
		assert.ErrorIs(t, err, classification.ErrUnmatched)
	})

	t.Run("webhook delete", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"delete": true,
			})
		}))
		defer server.Close()

		source, err := yamlSourceProvider{rawSourceProvider: rawTestSourceProvider{raw: []byte(`
workflows:
  default:
    - webhook:
        url: "` + server.URL + `"
`)}}.source()
		if err != nil {
			t.Fatal(err)
		}

		mocks := newTestClassifierMocks(t)
		workflow, err := mocks.compiler.Compile(source)
		if err != nil {
			t.Fatal(err)
		}

		_, err = workflow.Run(context.Background(), "default", Flags{}, model.Torrent{
			Name:        "Banned Thing.mkv",
			FilesStatus: model.FilesStatusSingle,
			Extension:   model.NewNullString("mkv"),
			Size:        1000000000,
		})
		assert.ErrorIs(t, err, classification.ErrDeleteTorrent)
	})

	t.Run("webhook non-2xx status is an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		}))
		defer server.Close()

		source, err := yamlSourceProvider{rawSourceProvider: rawTestSourceProvider{raw: []byte(`
workflows:
  default:
    - webhook:
        url: "` + server.URL + `"
`)}}.source()
		if err != nil {
			t.Fatal(err)
		}

		mocks := newTestClassifierMocks(t)
		workflow, err := mocks.compiler.Compile(source)
		if err != nil {
			t.Fatal(err)
		}

		_, err = workflow.Run(context.Background(), "default", Flags{}, model.Torrent{
			Name:        "Some Thing.mkv",
			FilesStatus: model.FilesStatusSingle,
			Extension:   model.NewNullString("mkv"),
			Size:        1000000000,
		})
		assert.Error(t, err)
	})
}
