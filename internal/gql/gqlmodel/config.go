package gqlmodel

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type ConfigQuery struct {
	Manager *config.Manager
}

func (q ConfigQuery) Params() []config.Param {
	return q.Manager.Params()
}

func (q ConfigQuery) Pending() bool {
	return q.Manager.HasPending()
}

type ConfigMutation struct {
	Manager *config.Manager
}

func (m ConfigMutation) Save(ref ref.Ref, value json_schema.JSONValue) (config.Param, error) {
	return m.Manager.Save(ref, value)
}

func (m ConfigMutation) Delete(ref ref.Ref) (config.Param, error) {
	return m.Manager.Delete(ref)
}
