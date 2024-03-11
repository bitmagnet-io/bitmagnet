package workflow

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestCompile(t *testing.T) {
	c := compiler{
		celEnvOption,
		conditions(
			andCondition{},
			//fileExtensionRatioCondition{},
			//fileTypeStatsCondition{},
			//hasContentTypeCondition{},
			//hasFilesStatusCondition{},
			//includesKeywords{},
			orCondition{},
			expressionCondition{},
		),
		actions(
			deleteAction{},
			findMatchAction{},
			ifElseAction{},
			noMatchAction{},
			noopAction{},
			parseVideoContentAction{},
			setContentTypeAction{},
			//sequenceAction{},
			//withVarsAction{},
		),
	}
	rawWorkflow := make(map[string]interface{})
	parseErr := yaml.Unmarshal([]byte(workflowDefaultYaml), &rawWorkflow)
	if parseErr != nil {
		t.Error(parseErr)
		return
	}
	t.Run("decode", func(t *testing.T) {
		workflow, compileErr := c.Compile(
			rawWorkflow,
		)
		if compileErr != nil {
			t.Fatal(compileErr)
			return
		}
		type testData struct {
			torrent     model.Torrent
			expected    Classification
			expectedErr error
		}
		testCases := []testData{
			{
				torrent: model.Torrent{
					Name:        "The XXX Movie 1080p.mkv",
					FilesStatus: model.FilesStatusSingle,
					Extension:   model.NewNullString("mkv"),
					Size:        1000000000,
					//Hint: model.TorrentHint{
					//	ContentType: model.ContentTypeXxx,
					//},
				},
				expected: Classification{
					ContentAttributes: classifier.ContentAttributes{
						ContentType:     model.NewNullContentType(model.ContentTypeXxx),
						VideoResolution: model.NewNullVideoResolution(model.VideoResolutionV1080p),
					},
				},
			},
			//{
			//	torrent: model.Torrent{
			//		name:        "The Regular Movie.mkv",
			//		FilesStatus: model.FilesStatusSingle,
			//		Extension:   model.NewNullString("mkv"),
			//		Size:        1000000000,
			//	},
			//	expected: Classification{
			//		ContentType: model.NewNullContentType(model.ContentTypeMovie),
			//	},
			//},
		}
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("torrent: %v", tc.torrent.Name), func(t *testing.T) {
				result, runErr := workflow.Run(context.Background(), tc.torrent)
				if runErr != nil {
					assert.Equal(t, tc.expectedErr, runErr)
					t.Log(runErr)
				} else {
					assert.Equal(t, tc.expected, result)
				}
			})
		}
	})
}
