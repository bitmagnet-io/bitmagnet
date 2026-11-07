package classifier

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

// matchlabTestCompleteResponse returns the complete four-field API contract.
const matchlabTestEpisodeID = "tt1234567:1:2"

func matchlabTestCompleteResponse(start, end string) string {
	return fmt.Sprintf(
		`{"sid_start":%q,"sid_end":%q,"asid_start":null,"asid_end":null}`,
		start, end,
	)
}

func TestMatchlabClient(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/match" || r.URL.Query().Get("title") != "Example.Show.S01E02" {
			t.Errorf("unexpected MatchLab request: %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(matchlabTestEpisodeID, matchlabTestEpisodeID))
	}))
	defer server.Close()
	client, err := newTestMatchlabClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	matched, found, err := client.Match(context.Background(), "Example.Show.S01E02")
	if err != nil || !found || matched.SIDStart == nil || *matched.SIDStart != matchlabTestEpisodeID {
		t.Fatalf("got %+v, %v, %v", matched, found, err)
	}
}

func TestMatchlabUnmatched(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"sid_start":null,"sid_end":null,"asid_start":null,"asid_end":null}`)
	}))
	defer server.Close()
	client, _ := newTestMatchlabClient(server.URL)
	_, found, err := client.Match(context.Background(), "Unknown")
	if err != nil || found {
		t.Fatalf("expected no match: %v, %v", found, err)
	}
}

func TestMatchlabRangeAndInvalidConfig(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, matchlabTestCompleteResponse("tt1234567:1:1", "tt1234567:1:5"))
	}))
	defer server.Close()
	client, _ := newTestMatchlabClient(server.URL)
	_, found, err := client.Match(context.Background(), "Pack")
	if err != nil || !found {
		t.Fatalf("expected IMDb identity from same-title range: %v, %v", found, err)
	}
	if _, err := newTestMatchlabClient("ftp://example.org"); err == nil {
		t.Fatal("expected invalid scheme")
	}
	if client, err := newTestMatchlabClient(""); err != nil || client != nil {
		t.Fatal("empty URL must disable MatchLab")
	}
}

func TestMatchlabParserOverridesExistingFields(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/parse" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, `{"title":"Example Show","year":2024,"seasons":[2],"episodes":[3,4]}`)
	}))
	defer server.Close()
	client, err := newTestMatchlabClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	cl := classification.Result{}
	cl.BaseTitle = model.NewNullString("Incorrect Title")
	cl.Episodes = model.Episodes{1: map[int]struct{}{8: {}}}
	if err := overrideMatchlabParsed(context.Background(), client, "Example.Show.S02E03-E04", &cl); err != nil {
		t.Fatal(err)
	}
	if cl.BaseTitle.String != "Example Show" || cl.Date.Year != 2024 || !cl.Episodes.HasEpisode(1, 8) {
		t.Fatalf("parser overrides not applied: %+v", cl)
	}
}

func TestMatchlabMatchOverridesEpisodes(t *testing.T) {
	t.Parallel()

	start, end := "tt1234567:2:3", "tt1234567:2:4"
	cl := classification.Result{ContentAttributes: classification.ContentAttributes{
		Episodes: model.Episodes{1: map[int]struct{}{8: {}}},
	}}
	applyMatchlabEpisodes(matchlab.Response{SIDStart: &start, SIDEnd: &end}, &cl)
	if !cl.Episodes.HasEpisode(2, 3) || !cl.Episodes.HasEpisode(2, 4) || cl.Episodes.HasEpisode(1, 8) {
		t.Fatalf("expected MatchLab SID episode override: %v", cl.Episodes)
	}
}

func TestMatchlabSeasonOverridesParsedEpisodes(t *testing.T) {
	t.Parallel()

	start, end := "tt1234567:2", "tt1234567:2"
	cl := classification.Result{ContentAttributes: classification.ContentAttributes{
		Episodes: model.Episodes{1: map[int]struct{}{8: {}}},
	}}

	applyMatchlabEpisodes(matchlab.Response{SIDStart: &start, SIDEnd: &end}, &cl)
	if len(cl.Episodes) != 1 || !cl.Episodes.HasEpisode(2, 8) || cl.Episodes.HasEpisode(1, 8) {
		t.Fatalf("expected MatchLab season to replace parsed episodes: %v", cl.Episodes)
	}
}

func TestMatchlabClassificationFromCoordinates(t *testing.T) {
	t.Parallel()

	ptr := func(v string) *string { return &v }
	cases := []struct {
		name     string
		response matchlab.Response
		want     model.ContentType
		found    bool
	}{
		{
			name:     "episode",
			response: matchlab.Response{SIDStart: ptr("tt1234567:2:3"), SIDEnd: ptr("tt1234567:2:4")},
			want:     model.ContentTypeTvShow,
			found:    true,
		},
		{
			name:     "season",
			response: matchlab.Response{SIDStart: ptr("tt1234567:2"), SIDEnd: ptr("tt1234567:2")},
			want:     model.ContentTypeTvShow,
			found:    true,
		},
		{
			name: "anime",
			response: matchlab.Response{
				SIDStart: ptr("tt1234567"), SIDEnd: ptr("tt1234567"), ASIDStart: ptr("1234:1"),
			},
			want:  model.ContentTypeTvShow,
			found: true,
		},
		{
			name:     "anime only",
			response: matchlab.Response{ASIDStart: ptr("1234:1"), ASIDEnd: ptr("1234:1")},
			want:     model.ContentTypeTvShow,
			found:    true,
		},
		{
			name:     "film",
			response: matchlab.Response{SIDStart: ptr("tt1234567"), SIDEnd: ptr("tt1234567")},
			want:     model.ContentTypeMovie,
			found:    true,
		},
		{
			name:     "empty",
			response: matchlab.Response{},
			want:     "",
			found:    false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, ok := matchlabContentType(tc.response)
			if got != tc.want || ok != tc.found {
				t.Fatalf("got %q %v, want %q %v", got, ok, tc.want, tc.found)
			}
		})
	}
}

func newTestMatchlabClient(rawURL string) (matchlab.Client, error) {
	return matchlab.NewClient(matchlab.Config{URL: rawURL, TimeoutSeconds: 8, MaxConcurrent: 4})
}
