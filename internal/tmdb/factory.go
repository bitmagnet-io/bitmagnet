package tmdb

import (
	"go.uber.org/zap"
)

func New(config Config, logger *zap.SugaredLogger) Client {
	return client{
		requester: &requesterLazy{
			config: config,
			logger: logger.Named(Namespace),
		},
	}
}
