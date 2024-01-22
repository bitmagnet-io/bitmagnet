package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sort"
)

type Params struct {
	fx.In
	SubClassifiers []lazy.Lazy[SubClassifier] `group:"content_classifiers"`
	Logger         *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Classifier lazy.Lazy[Classifier]
}

func New(p Params) Result {
	return Result{
		Classifier: lazy.New(func() (Classifier, error) {
			subClassifiers := make([]SubClassifier, 0, len(p.SubClassifiers)+1)
			for _, subResolver := range p.SubClassifiers {
				r, err := subResolver.Get()
				if err != nil {
					return nil, err
				}
				subClassifiers = append(subClassifiers, r)
			}
			subClassifiers = append(subClassifiers, fallbackClassifier{})
			sort.Slice(subClassifiers, func(i, j int) bool {
				return subClassifiers[i].Priority() < subClassifiers[j].Priority()
			})
			return classifier{subClassifiers, p.Logger}, nil
		}),
	}
}
