package classifier

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"math"
)

type FallbackClassifier struct{}

func (c FallbackClassifier) Key() string {
	return "fallback"
}

func (c FallbackClassifier) Priority() int {
	return math.MaxInt
}

func (c FallbackClassifier) Classify(_ context.Context, t model.Torrent) (Classification, error) {
	cl := Classification{}
	cl.ApplyHint(t.Hint)
	return cl, nil
}
