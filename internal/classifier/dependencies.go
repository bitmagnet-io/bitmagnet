package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"go.uber.org/zap"
)

type dependencies struct {
	search         LocalSearch
	tmdbClient     tmdb.Client
	matchlabClient matchlab.Client
	tmdbEnabled    bool
	logger         *zap.SugaredLogger
}
