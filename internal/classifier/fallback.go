package classifier

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"math"
)

type fallbackClassifier struct{}

func (c fallbackClassifier) Key() string {
	return "fallback"
}

func (c fallbackClassifier) Priority() int {
	return math.MaxInt
}

func (c fallbackClassifier) Classify(_ context.Context, t model.Torrent) (Classification, error) {
	cl := Classification{}
	cl.ApplyHint(t.Hint)
	return cl, nil
}
