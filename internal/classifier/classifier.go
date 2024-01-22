package classifier

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"go.uber.org/zap"
)

var (
	ErrNoMatch = errors.New("no match")
)

type Classifier interface {
	Classify(ctx context.Context, torrent model.Torrent) (Classification, error)
}

type SubClassifier interface {
	Classifier
	Key() string
	Priority() int
}

type classifier struct {
	subClassifiers []SubClassifier
	logger         *zap.SugaredLogger
}

func (c classifier) Classify(ctx context.Context, t model.Torrent) (Classification, error) {
	for _, sc := range c.subClassifiers {
		tc, err := sc.Classify(ctx, t)
		if err == nil {
			return tc, nil
		}
		if !errors.Is(err, ErrNoMatch) {
			c.logger.Errorw("error classifying content", "classifier", sc.Key(), "torrent", t, "error", err)
			return Classification{}, err
		}
	}
	return Classification{}, ErrNoMatch
}
