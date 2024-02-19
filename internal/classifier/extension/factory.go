package extension

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Classifier lazy.Lazy[classifier.SubClassifier] `group:"content_classifiers"`
}

func New() Result {
	return Result{
		Classifier: lazy.New(func() (classifier.SubClassifier, error) {
			return extensionClassifier{}, nil
		}),
	}
}
