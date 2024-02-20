package keywords

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Xxx lazy.Lazy[classifier.SubClassifier] `group:"content_classifiers"`
}

func New() Result {
	return Result{
		Xxx: lazy.New(func() (classifier.SubClassifier, error) {
			return NewKeywordsClassifier(
				model.ContentTypeXxx,
				xxxWords,
				20,
			), nil
		}),
	}
}
