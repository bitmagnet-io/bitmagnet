package tmdb

import (
	"go.uber.org/zap"
)

func New(config Config, logger *zap.Logger) Client {
	return client{
		requester: &requesterLazy{
			config: config,
			logger: logger,
		},
	}
}
