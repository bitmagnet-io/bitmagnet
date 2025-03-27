package classifier

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	classifier_mocks "github.com/bitmagnet-io/bitmagnet/internal/classifier/mocks"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	tmdb_mocks "github.com/bitmagnet-io/bitmagnet/internal/tmdb/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClassifier(t *testing.T) {
	t.Parallel()

	matchContext := mock.MatchedBy(func(ctx any) bool {
		_, ok := ctx.(context.Context)
		return ok
	})

	testCases := []struct {
		torrent      model.Torrent
		flags        Flags
		prepareMocks func(mocks testClassifierMocks)
		expected     classification.Result
		expectedErr  error
	}{
		{
			torrent: model.Torrent{
				Name:        "The Regular Movie (2000).mkv",
				FilesStatus: model.FilesStatusSingle,
				Extension:   model.NewNullString("mkv"),
				Size:        1000000000,
			},
			prepareMocks: func(mocks testClassifierMocks) {
				mocks.search.On(
					"ContentBySearch",
					matchContext,
					model.ContentTypeMovie,
					"The Regular Movie",
					model.Year(2000),
				).
					Return(model.Content{}, classification.ErrUnmatched)
				mocks.tmdbClient.On(
					"SearchMovie",
					matchContext,
					tmdb.SearchMovieRequest{
						Query:        "The Regular Movie",
						Year:         2000,
						IncludeAdult: true,
					},
				).
					Return(tmdb.SearchMovieResponse{}, nil)
			},
			expected: classification.Result{
				ContentAttributes: classification.ContentAttributes{
					ContentType: model.NewNullContentType(model.ContentTypeMovie),
					BaseTitle:   model.NewNullString("The Regular Movie"),
					Date: model.Date{
						Year: 2000,
					},
				},
			},
		},
		{
			torrent: model.Torrent{
				Name:        "The Regular Local Movie (2000).mkv",
				FilesStatus: model.FilesStatusSingle,
				Extension:   model.NewNullString("mkv"),
				Size:        1000000000,
			},
			prepareMocks: func(mocks testClassifierMocks) {
				mocks.search.On(
					"ContentBySearch",
					matchContext,
					model.ContentTypeMovie,
					"The Regular Local Movie",
					model.Year(2000),
				).
					Return(model.Content{
						Type:        model.ContentTypeMovie,
						Source:      "local",
						ID:          "123",
						Title:       "The Regular Local Movie",
						ReleaseYear: 2000,
					}, nil)
			},
			expected: classification.Result{
				ContentAttributes: classification.ContentAttributes{
					ContentType: model.NewNullContentType(model.ContentTypeMovie),
					BaseTitle:   model.NewNullString("The Regular Local Movie"),
					Date: model.Date{
						Year: 2000,
					},
				},
				Content: &model.Content{
					Type:        model.ContentTypeMovie,
					Source:      "local",
					ID:          "123",
					Title:       "The Regular Local Movie",
					ReleaseYear: 2000,
				},
			},
		},
		{
			torrent: model.Torrent{
				Name:        "The Regular TMDB Movie (2000).mkv",
				FilesStatus: model.FilesStatusSingle,
				Extension:   model.NewNullString("mkv"),
				Size:        1000000000,
			},
			prepareMocks: func(mocks testClassifierMocks) {
				mocks.search.On(
					"ContentBySearch",
					matchContext,
					model.ContentTypeMovie,
					"The Regular TMDB Movie",
					model.Year(2000),
				).
					Return(model.Content{}, classification.ErrUnmatched)
				mocks.tmdbClient.On(
					"SearchMovie",
					matchContext,
					tmdb.SearchMovieRequest{
						Query:        "The Regular TMDB Movie",
						Year:         2000,
						IncludeAdult: true,
					},
				).
					Return(tmdb.SearchMovieResponse{
						Results: []tmdb.SearchMovieResult{
							{
								ID:          123,
								Title:       "The Regular TMDB Movie",
								ReleaseDate: "2000-01-01",
							},
						},
					}, nil)
				mocks.tmdbClient.On(
					"MovieDetails",
					matchContext,
					tmdb.MovieDetailsRequest{
						ID: 123,
					},
				).
					Return(tmdb.MovieDetailsResponse{
						ID:            123,
						Title:         "The Regular TMDB Movie",
						OriginalTitle: "The Regular TMDB Movie Original",
						ReleaseDate:   "2000-01-01",
					}, nil)
			},
			expected: classification.Result{
				ContentAttributes: classification.ContentAttributes{
					ContentType: model.NewNullContentType(model.ContentTypeMovie),
					BaseTitle:   model.NewNullString("The Regular TMDB Movie"),
					Date: model.Date{
						Year: 2000,
					},
				},
				Content: &model.Content{
					Type:   model.ContentTypeMovie,
					Source: "tmdb",
					ID:     "123",
					Title:  "The Regular TMDB Movie",
					ReleaseDate: model.Date{
						Year:  2000,
						Month: 1,
						Day:   1,
					},
					ReleaseYear:   2000,
					Adult:         model.NewNullBool(false),
					OriginalTitle: model.NewNullString("The Regular TMDB Movie Original"),
					Popularity:    model.NewNullFloat32(0),
					VoteAverage:   model.NewNullFloat32(0),
					VoteCount:     model.NewNullUint(0),
				},
			},
		},
		{
			torrent: model.Torrent{
				Name:        "The XXX Movie 1080p.mkv",
				FilesStatus: model.FilesStatusSingle,
				Extension:   model.NewNullString("mkv"),
				Size:        1000000000,
			},
			expected: classification.Result{
				ContentAttributes: classification.ContentAttributes{
					ContentType:     model.NewNullContentType(model.ContentTypeXxx),
					VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("torrent: %s", tc.torrent.Name), func(t *testing.T) {
			t.Parallel()

			mocks := newTestClassifierMocks(t)
			source, sourceErr := coreSourceProvider{}.provider().source()

			if sourceErr != nil {
				t.Fatal(sourceErr)
				return
			}

			workflow, compileErr := mocks.compiler.Compile(source)
			if compileErr != nil {
				t.Fatal(compileErr)
				return
			}

			if tc.prepareMocks != nil {
				tc.prepareMocks(mocks)
			}

			result, runErr := workflow.Run(context.Background(), "default", tc.flags, tc.torrent)
			if runErr != nil {
				assert.Equal(t, tc.expectedErr, runErr)
				t.Log(runErr)
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

type testClassifierMocks struct {
	compiler   Compiler
	search     *classifier_mocks.LocalSearch
	tmdbClient *tmdb_mocks.Client
}

func newTestClassifierMocks(t *testing.T) testClassifierMocks {
	t.Helper()

	search := classifier_mocks.NewLocalSearch(t)
	tmdbClient := tmdb_mocks.NewClient(t)

	return testClassifierMocks{
		compiler: compiler{
			options: []compilerOption{
				compilerFeatures(defaultFeatures),
				celEnvOption,
			},
			dependencies: dependencies{
				search:     search,
				tmdbClient: tmdbClient,
			},
		},
		search:     search,
		tmdbClient: tmdbClient,
	}
}
