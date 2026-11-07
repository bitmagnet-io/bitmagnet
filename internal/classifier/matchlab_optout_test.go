package classifier

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	tmdb_mocks "github.com/bitmagnet-io/bitmagnet/internal/tmdb/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	optOutMatchLabFlag    = "matchlab_enabled"
	optOutLocalSearchFlag = "local_search_enabled"
	optOutAPIsFlag        = "apis_enabled"
	optOutLocalSource     = "local"
	optOutLocalTitle      = "The Regular Local Movie"
	optOutIMDbID          = "tt1234567"
)

func TestNativeSearchAfterMatchLabDisconnectWithStaleQueuedFlag(t *testing.T) {
	t.Parallel()

	mocks := newTestClassifierMocks(t)
	source, err := (yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}).source()
	require.NoError(t, err)
	workflow, err := mocks.compiler.Compile(source)
	require.NoError(t, err)

	torrent := model.Torrent{
		Name:        "The Regular Local Movie (2000).mkv",
		FilesStatus: model.FilesStatusSingle,
		Extension:   model.NewNullString("mkv"),
		Size:        1000000000,
	}
	content := model.Content{
		Type: model.ContentTypeMovie, Source: optOutLocalSource, ID: "42",
		Title: optOutLocalTitle, ReleaseYear: 2000,
	}
	mocks.search.On("ContentBySearch", mock.Anything, model.ContentTypeMovie,
		optOutLocalTitle, model.Year(2000)).Return(content, nil).Once()

	// A process-torrent job persisted while MatchLab was configured may
	// contain matchlab_enabled=true even after its URL has been removed.
	// Native search must work without MatchLab or TMDB network access.
	result, err := workflow.Run(context.Background(), "default", Flags{
		optOutMatchLabFlag:    true,
		optOutLocalSearchFlag: true,
		optOutAPIsFlag:        false,
	}, torrent)
	require.NoError(t, err)
	require.NotNil(t, result.Content)
	require.Equal(t, optOutLocalSource, result.Content.Source)
	require.Equal(t, "42", result.Content.ID)
}

func TestNativeModePreservesPreviouslyStoredMatchLabIMDbContent(t *testing.T) {
	t.Parallel()

	mocks := newTestClassifierMocks(t)
	source, err := (yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}).source()
	require.NoError(t, err)
	workflow, err := mocks.compiler.Compile(source)
	require.NoError(t, err)

	stored := model.Content{
		Type:   model.ContentTypeMovie,
		Source: model.SourceImdb,
		ID:     optOutIMDbID,
		Title:  "Verified Movie",
	}
	torrent := model.Torrent{
		Name: "Verified.Movie.2024.mkv",
		Hint: model.TorrentHint{
			ContentType:   model.ContentTypeMovie,
			ContentSource: model.NewNullString(model.SourceImdb),
			ContentID:     model.NewNullString(optOutIMDbID),
		},
		Contents: []model.TorrentContent{{
			ContentType:   model.NewNullContentType(model.ContentTypeMovie),
			ContentSource: model.NewNullString(model.SourceImdb),
			ContentID:     model.NewNullString(optOutIMDbID),
			Content:       stored,
		}},
	}
	// The normal Bitmagnet classification path reuses an existing verified
	// hint. Disconnecting MatchLab must not erase usable stored metadata.
	result, err := workflow.Run(context.Background(), "default", Flags{
		optOutMatchLabFlag:    true,
		optOutLocalSearchFlag: true,
		optOutAPIsFlag:        false,
	}, torrent)
	require.NoError(t, err)
	require.NotNil(t, result.Content)
	require.Equal(t, stored.Source, result.Content.Source)
	require.Equal(t, stored.ID, result.Content.ID)
}

func TestDisconnectedMatchLabDoesNotForceNativeSearchAgainstUserChoice(t *testing.T) {
	t.Parallel()

	mocks := newTestClassifierMocks(t)
	source, err := (yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}).source()
	require.NoError(t, err)
	workflow, err := mocks.compiler.Compile(source)
	require.NoError(t, err)

	torrent := model.Torrent{
		Name:        "The Regular Movie (2000).mkv",
		FilesStatus: model.FilesStatusSingle,
		Extension:   model.NewNullString("mkv"),
		Size:        1000000000,
	}
	result, err := workflow.Run(context.Background(), "default", Flags{
		optOutMatchLabFlag:    true,
		optOutLocalSearchFlag: false,
		optOutAPIsFlag:        false,
	}, torrent)
	require.NoError(t, err)
	require.Nil(t, result.Content)
}

func TestConnectThenDisconnectMatchLabRetainsContentAndNativeMode(t *testing.T) {
	t.Parallel()

	api := testMatchLab(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case matchlabTestMatchPath:
			_, _ = fmt.Fprint(w, matchlabTestCompleteResponse(optOutIMDbID, optOutIMDbID))
		case matchlabTestParsePath:
			_, _ = fmt.Fprint(w, `{"title":"Verified Movie","year":2024}`)
		default:
			t.Errorf("unexpected MatchLab request %q", r.URL.Path)
		}
	})
	remote, torrent := runWithMatchLab(t, api, false, tmdb_mocks.NewClient(t))
	identified, err := remote.Run(context.Background(), "default", nil, torrent)
	require.NoError(t, err)
	require.NotNil(t, identified.Content)
	require.Equal(t, model.SourceImdb, identified.Content.Source)
	require.Equal(t, optOutIMDbID, identified.Content.ID)

	// Simulate the saved torrent_contents record becoming a content hint
	// during normal reprocessing, then restart without MATCHLAB_URL.
	stored := *identified.Content
	torrent.Hint = model.TorrentHint{
		ContentType:   stored.Type,
		ContentSource: model.NewNullString(stored.Source),
		ContentID:     model.NewNullString(stored.ID),
	}
	torrent.Contents = []model.TorrentContent{{
		ContentType:   model.NewNullContentType(stored.Type),
		ContentSource: model.NewNullString(stored.Source),
		ContentID:     model.NewNullString(stored.ID),
		Content:       stored,
	}}
	nativeMocks := newTestClassifierMocks(t)
	source, err := (yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}).source()
	require.NoError(t, err)
	native, err := nativeMocks.compiler.Compile(source)
	require.NoError(t, err)
	restored, err := native.Run(context.Background(), "default", Flags{
		optOutMatchLabFlag:    true, // persisted before disconnection
		optOutLocalSearchFlag: true,
		optOutAPIsFlag:        false,
	}, torrent)
	require.NoError(t, err)
	require.NotNil(t, restored.Content)
	require.Equal(t, stored.ID, restored.Content.ID)
	require.Equal(t, stored.Source, restored.Content.Source)
}
