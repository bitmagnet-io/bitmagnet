package keywords

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Music     lazy.Lazy[classifier.SubClassifier] `group:"content_classifiers"`
	Audiobook lazy.Lazy[classifier.SubClassifier] `group:"content_classifiers"`
	Ebook     lazy.Lazy[classifier.SubClassifier] `group:"content_classifiers"`
	Xxx       lazy.Lazy[classifier.SubClassifier] `group:"content_classifiers"`
}

func New() Result {
	return Result{
		Audiobook: lazy.New(func() (classifier.SubClassifier, error) {
			return keywordsClassifier{
				contentType:       model.ContentTypeAudiobook,
				regex:             regex.NewRegexFromNames(audiobookWords...),
				priority:          20,
				requiredFileTypes: []model.FileType{model.FileTypeAudio},
			}, nil
		}),
		Music: lazy.New(func() (classifier.SubClassifier, error) {
			return keywordsClassifier{
				contentType:       model.ContentTypeMusic,
				regex:             regex.NewRegexFromNames(musicWords...),
				priority:          21,
				requiredFileTypes: []model.FileType{model.FileTypeAudio},
			}, nil
		}),
		Ebook: lazy.New(func() (classifier.SubClassifier, error) {
			return keywordsClassifier{
				contentType:       model.ContentTypeAudiobook,
				regex:             regex.NewRegexFromNames(ebookWords...),
				priority:          22,
				requiredFileTypes: []model.FileType{model.FileTypeDocument},
			}, nil
		}),
		Xxx: lazy.New(func() (classifier.SubClassifier, error) {
			return keywordsClassifier{
				contentType: model.ContentTypeXxx,
				regex:       regex.NewRegexFromNames(xxxWords...),
				priority:    23,
			}, nil
		}),
	}
}
