package classifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	classifier_mocks "github.com/bitmagnet-io/bitmagnet/internal/classifier/mocks"
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	tmdb_mocks "github.com/bitmagnet-io/bitmagnet/internal/tmdb/mocks"
	"github.com/stretchr/testify/mock"
)

func runWithMatchLab(t *testing.T, api matchlab.Client, enabledTMDB bool, client tmdb.Client) (Runner, model.Torrent) {
	t.Helper()
	source, err := (yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}).source()
	if err != nil {
		t.Fatal(err)
	}
	search := classifier_mocks.NewLocalSearch(t)
	compiler := compiler{
		options: []compilerOption{compilerFeatures(defaultFeatures), celEnvOption},
		dependencies: dependencies{
			search:         search,
			tmdbClient:     client,
			matchlabClient: api,
			tmdbEnabled:    enabledTMDB,
		},
	}
	workflow, err := compiler.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	return workflow, model.Torrent{
		Name:        "Example.Movie.2024.1080p.WEB-DL.mkv",
		FilesStatus: model.FilesStatusSingle,
		Extension:   model.NewNullString("mkv"),
		Size:        1_000_000_000,
	}
}

func testMatchLab(t *testing.T, handler http.HandlerFunc) matchlab.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	api, err := matchlab.NewClient(matchlab.Config{URL: server.URL, TimeoutSeconds: 4, MaxConcurrent: 2})
	if err != nil {
		t.Fatal(err)
	}
	return api
}

const (
	matchlabTestMatchPath = "/match"
	matchlabTestParsePath = "/parse"
	matchlabTestIMDbID    = "tt1234567"
)

const matchlabTestExternalSource = "imdb_id"

func TestConfiguredMatchLabDisablesAllLocalAndTMDBTitleSearch(t *testing.T) {
	t.Parallel()
	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case matchlabTestMatchPath:
			_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(matchlabTestIMDbID, matchlabTestIMDbID))
		case matchlabTestParsePath:
			_, _ = fmt.Fprint(w, `{"title":"Example Movie","year":"2024"}`)
		default:
			t.Errorf("unexpected MatchLab path %q", r.URL.Path)
		}
	})
	tmdbClient := tmdb_mocks.NewClient(t)
	workflow, torrent := runWithMatchLab(t, api, false, tmdbClient)
	// Even an old reprocessing job explicitly requesting local matching must
	// not bypass a configured MatchLab provider.
	result, err := workflow.Run(context.Background(), "default", Flags{
		"matchlab_enabled": false, "local_search_enabled": true,
	}, torrent)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil || result.Content.Source != model.SourceImdb ||
		result.Content.ID != matchlabTestIMDbID || result.Content.Title != "Example Movie" {
		t.Fatalf("unexpected verified IMDb content: %+v", result.Content)
	}
}

func TestMatchLabSkipsKnownNonVideoTorrentCategories(t *testing.T) {
	t.Parallel()
	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("non-video torrent must not trigger MatchLab %s", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	workflow, torrent := runWithMatchLab(t, api, false, tmdb_mocks.NewClient(t))
	torrent.Name = "Accounting.Application.v2.1.zip"
	torrent.Extension = model.NewNullString("zip")
	torrent.Size = 5_000_000
	torrent.Hint = model.TorrentHint{ContentType: model.ContentTypeSoftware}

	result, err := workflow.Run(context.Background(), "default", nil, torrent)
	if err != nil {
		t.Fatal(err)
	}
	if result.ContentType.ContentType != model.ContentTypeSoftware || result.Content != nil {
		t.Fatalf("MatchLab reclassified software torrent: %+v", result)
	}
}

func TestMatchLabNoMatchAndOutageNeverFallback(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		code      int
		payload   string
		wantError bool
	}{
		{"unmatched", 200, `{"sid_start":null,"sid_end":null,"asid_start":null,"asid_end":null}`, false},
		{"unavailable", 503, `{"error":"initializing"}`, true},
		{"malformed", 200, matchlabTestCompleteResponse(matchlabTestIMDbID, "tt2345678"), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != matchlabTestMatchPath {
					t.Errorf("unexpected request %s", r.URL.Path)
				}
				w.WriteHeader(tc.code)
				_, _ = fmt.Fprint(w, tc.payload)
			})
			workflow, torrent := runWithMatchLab(t, api, false, tmdb_mocks.NewClient(t))
			result, err := workflow.Run(context.Background(), "default", nil, torrent)
			if (err != nil) != tc.wantError {
				t.Fatalf("error=%v, content=%+v", err, result.Content)
			}
			if errors.Is(err, classification.ErrUnmatched) {
				t.Fatalf("outage should not be a legitimate no-match: %v", err)
			}
			if result.Content != nil {
				t.Fatalf("unexpected fallback content: %+v", result.Content)
			}
		})
	}
}

func TestMatchLabIMDbIdentitySurvivesTMDBEnrichment(t *testing.T) {
	t.Parallel()
	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case matchlabTestMatchPath:
			_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(matchlabTestIMDbID, matchlabTestIMDbID))
		case matchlabTestParsePath:
			t.Errorf("redundant /parse request after successful TMDB metadata lookup")
		default:
			t.Errorf("unexpected API path %s", r.URL.Path)
		}
	})
	client := tmdb_mocks.NewClient(t)
	var lookup tmdb.FindByIDResponse
	if err := json.Unmarshal([]byte(`{"movie_results":[{"id":42}]}`), &lookup); err != nil {
		t.Fatal(err)
	}
	matchContext := mock.MatchedBy(func(ctx any) bool { _, ok := ctx.(context.Context); return ok })
	client.On("FindByID", matchContext, tmdb.FindByIDRequest{
		ExternalSource: matchlabTestExternalSource, ExternalID: matchlabTestIMDbID,
	}).Return(lookup, nil).Once()
	client.On("MovieDetails", matchContext, tmdb.MovieDetailsRequest{ID: 42}).
		Return(tmdb.MovieDetailsResponse{
			ID: 42, Title: "TMDB Movie",
			IMDbID: matchlabTestIMDbID, Overview: "Enriched metadata",
		}, nil).Once()
	workflow, torrent := runWithMatchLab(t, api, true, client)
	result, err := workflow.Run(context.Background(), "default", nil, torrent)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil || result.Content.ID != matchlabTestIMDbID ||
		result.Content.Source != model.SourceImdb || result.Content.Title != "TMDB Movie" {
		t.Fatalf("TMDB changed canonical MatchLab identity: %+v", result.Content)
	}
	hasTMDBLink := false
	for _, attr := range result.Content.Attributes {
		if attr.Source == model.SourceTmdb && attr.Key == "id" && attr.Value == "42" {
			hasTMDBLink = true
		}
	}
	if !hasTMDBLink {
		t.Fatalf("TMDB external link missing: %+v", result.Content.Attributes)
	}
}

func TestMatchLabBareIMDbIDUsesTMDBMediaTypeWithoutTitleSearch(t *testing.T) {
	t.Parallel()
	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case matchlabTestMatchPath:
			_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(matchlabTestIMDbID, matchlabTestIMDbID))
		case matchlabTestParsePath:
			_, _ = fmt.Fprint(w, `{"title":"Example Show","year":2024}`)
		default:
			t.Errorf("unexpected request %s", r.URL.Path)
		}
	})
	client := tmdb_mocks.NewClient(t)
	var lookup tmdb.FindByIDResponse
	if err := json.Unmarshal([]byte(`{"tv_results":[{"id":42}]}`), &lookup); err != nil {
		t.Fatal(err)
	}
	client.On("FindByID", mock.Anything, tmdb.FindByIDRequest{
		ExternalSource: matchlabTestExternalSource, ExternalID: matchlabTestIMDbID,
	}).Return(lookup, nil).Once()
	client.On("TvDetails", mock.Anything, tmdb.TvDetailsRequest{
		SeriesID: 42, AppendToResponse: []string{"external_ids"},
	}).Return(tmdb.TvDetailsResponse{ID: 42, Name: "Example Show"}, nil).Once()
	workflow, torrent := runWithMatchLab(t, api, true, client)
	result, err := workflow.Run(context.Background(), "default", nil, torrent)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil || result.Content.Type != model.ContentTypeTvShow ||
		result.Content.Source != model.SourceImdb || result.Content.ID != matchlabTestIMDbID {
		t.Fatalf("bare SID should use TV metadata type without changing IMDb identity: %+v", result.Content)
	}
}

func TestMatchLabKeepsIMDbIdentityWhenTMDBReturnsNotFound(t *testing.T) {
	t.Parallel()
	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case matchlabTestMatchPath:
			_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(matchlabTestIMDbID, matchlabTestIMDbID))
		case matchlabTestParsePath:
			_, _ = fmt.Fprint(w, `{"title":"Example Movie","year":2024}`)
		default:
			t.Errorf("unexpected request %s", r.URL.Path)
		}
	})
	client := tmdb_mocks.NewClient(t)
	client.On("FindByID", mock.Anything, tmdb.FindByIDRequest{
		ExternalSource: matchlabTestExternalSource, ExternalID: matchlabTestIMDbID,
	}).Return(tmdb.FindByIDResponse{}, tmdb.ErrNotFound).Once()
	workflow, torrent := runWithMatchLab(t, api, true, client)
	result, err := workflow.Run(context.Background(), "default", nil, torrent)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil || result.Content.Source != model.SourceImdb ||
		result.Content.ID != matchlabTestIMDbID || result.Content.Type != model.ContentTypeMovie {
		t.Fatalf("TMDB not-found must retain the verified IMDb identity: %+v", result.Content)
	}
}

func TestMatchLabOverridesStaleMovieHintAndPreviouslyAttachedContent(t *testing.T) {
	t.Parallel()
	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case matchlabTestMatchPath:
			_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(matchlabTestIMDbID, matchlabTestIMDbID))
		case matchlabTestParsePath:
			_, _ = fmt.Fprint(w, `{"title":"Example Movie","year":2024}`)
		default:
			t.Errorf("unexpected MatchLab path %q", r.URL.Path)
		}
	})
	workflow, torrent := runWithMatchLab(t, api, false, tmdb_mocks.NewClient(t))
	oldID := "tt7654321"
	torrent.Hint = model.TorrentHint{
		ContentType:   model.ContentTypeMovie,
		ContentSource: model.NewNullString(model.SourceImdb),
		ContentID:     model.NewNullString(oldID),
	}
	torrent.Contents = []model.TorrentContent{{
		ContentType:   model.NewNullContentType(model.ContentTypeMovie),
		ContentSource: model.NewNullString(model.SourceImdb),
		ContentID:     model.NewNullString(oldID),
		Content: model.Content{
			Type: model.ContentTypeMovie, Source: model.SourceImdb,
			ID: oldID, Title: "Stale IMDb identity",
		},
	}}
	result, err := workflow.Run(context.Background(), "default", nil, torrent)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil || result.Content.Source != model.SourceImdb ||
		result.Content.ID != matchlabTestIMDbID {
		t.Fatalf("MatchLab must replace stale cached identity: %+v", result.Content)
	}
}

func TestMatchLabBlocksDirectIdentitySearchActionsInCustomWorkflows(t *testing.T) {
	t.Parallel()
	api, err := matchlab.NewClient(matchlab.Config{
		URL: "http://127.0.0.1:1", TimeoutSeconds: 2, MaxConcurrent: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		def  actionDefinition
	}{
		{"local title search", attachLocalContentBySearchAction{}},
		{"TMDB title search", attachTmdbContentBySearchAction{}},
		{"local ID lookup", attachLocalContentByIDAction{}},
		{"TMDB ID lookup", attachTMDBContentByIDAction{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			compiled, compileErr := tc.def.compileAction(compilerContext{source: tc.def.name()})
			if compileErr != nil {
				t.Fatal(compileErr)
			}
			_, runErr := compiled.run(executionContext{
				Context:      context.Background(),
				dependencies: dependencies{matchlabClient: api},
				torrent: model.Torrent{Hint: model.TorrentHint{
					ContentType:   model.ContentTypeMovie,
					ContentSource: model.NewNullString(model.SourceImdb),
					ContentID:     model.NewNullString("tt7654321"),
				}},
				result: classification.Result{
					ContentAttributes: classification.ContentAttributes{
						BaseTitle:   model.NewNullString("Example Movie"),
						ContentType: model.NewNullContentType(model.ContentTypeMovie),
					},
				},
			})
			if !errors.Is(runErr, classification.ErrUnmatched) {
				t.Fatalf("direct search action bypassed configured MatchLab: %v", runErr)
			}
		})
	}
}
