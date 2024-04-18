package classifier

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	classifier_mocks "github.com/bitmagnet-io/bitmagnet/internal/classifier/mocks"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	tmdb_mocks "github.com/bitmagnet-io/bitmagnet/internal/tmdb/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestClassifier(t *testing.T) {
	matchContext := mock.MatchedBy(func(ctx any) bool {
		_, ok := ctx.(context.Context)
		return ok
	})
	testCases := []struct {
		torrent      model.Torrent
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
					Return(model.Content{}, classification.ErrNoMatch)
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
			mocks := newTestClassifierMocks(t)
			source, sourceErr := yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}.source()
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
			result, runErr := workflow.Run(context.Background(), "default", tc.torrent)
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
