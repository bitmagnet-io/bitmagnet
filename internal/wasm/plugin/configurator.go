package plugin

import (
	context "context"

	"github.com/bitmagnet-io/bitmagnet/proto/host/configurator"
)

type configuratorImpl struct {
	jsonConfig string
}

func (c configuratorImpl) GetConfig(context.Context, *configurator.Empty) (*configurator.Config, error) {
	return &configurator.Config{
		Json: c.jsonConfig,
	}, nil
}
