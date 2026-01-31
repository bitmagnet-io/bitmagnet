package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type Param struct {
	Ref        ref.Ref               `json:"ref"`
	Plugin     ref.Ref               `json:"plugin"`
	Value      json_schema.JSONValue `json:"value"`
	Source     string                `json:"source"`
	Default    json_schema.JSONValue `json:"default"`
	Dynamic    bool                  `json:"dynamic"`
	Pending    bool                  `json:"pending"`
	JSONSchema json_schema.JSONValue `json:"jsonSchema"`
}

func transformParam(param *resolver.Param) Param {
	defaultValue, _ := param.EncodeYAMLAny(param.NewDefaultAny())

	return Param{
		Ref:        param.Ref,
		Plugin:     param.Plugin,
		Value:      json_schema.JSONValue(param.ValueYAML()),
		Source:     param.Source(),
		Default:    json_schema.JSONValue(defaultValue),
		Dynamic:    param.IsDynamic(),
		Pending:    param.IsPending(),
		JSONSchema: json_schema.MustNewValue(param.JSONSchema()),
	}
}
