package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	"github.com/bitmagnet-io/plugin-opensearch/i18n"
	"github.com/bitmagnet-io/plugin-opensearch/shared"
)

type Config struct {
	Addresses        []string `json:"addresses"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	IndexPrefix      string   `json:"index_prefix"`
	NumberOfShards   int      `json:"number_of_shards"`
	NumberOfReplicas int      `json:"number_of_replicas"`
}

func (c Config) FullIndexPrefix() string {
	return fmt.Sprintf("%sv%d-", c.IndexPrefix, shared.TemplateVersion)
}

type param struct {
	description string
	schema      json_schema.JSONSchema
}

var params = map[string]param{
	"addresses": {
		description: "Addresses of OpenSearch nodes",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeArray),
			json_schema.Items(json_schema.MustNew(
				json_schema.Typed(json_schema.TypeString),
			)),
			json_schema.Default(json_schema.MustNewValue([]any{"http://localhost:9200"})),
		),
	},
	"username": {
		description: "Username",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
			json_schema.Default(json_schema.MustNewValue("admin")),
		),
	},
	"password": {
		description: "Password",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
			json_schema.Default(json_schema.MustNewValue("admin")),
		),
	},
	"index_prefix": {
		description: "Prefix for index names",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
			json_schema.Default(json_schema.MustNewValue("bitmagnet_")),
		),
	},
	"number_of_shards": {
		description: "Number of shards",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeInteger),
			json_schema.Default(json_schema.MustNewValue(1)),
		),
	},
	"number_of_replicas": {
		description: "Number of replicas",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeInteger),
			json_schema.Default(json_schema.MustNewValue(1)),
		),
	},
}

var Params = func() []*plugin.ConfigParam {
	result := make([]*plugin.ConfigParam, 0, len(params))

	for _, k := range slices.Sorted(maps.Keys(params)) {
		param := params[k]
		schema, err := json.Marshal(param.schema)
		if err != nil {
			panic(err)
		}

		result = append(result, &plugin.ConfigParam{
			Name:   k,
			Schema: schema,
		})
	}

	return result
}()

func LocalizedContent(localizer *i18n.Localizer) []*plugin.ConfigParamLocalizedContent {
	result := make([]*plugin.ConfigParamLocalizedContent, 0, len(params))
	for key := range params {
		result = append(result, &plugin.ConfigParamLocalizedContent{
			Name:        key,
			Description: localizer.Localize(key),
		})
	}
	return result
}

func I18NMessages() []*i18n.Message {
	messages := make([]*i18n.Message, 0, len(params))

	for _, k := range slices.Sorted(maps.Keys(params)) {
		param := params[k]
		messages = append(messages, &i18n.Message{
			ID:          k,
			Description: "Description for config field " + k,
			Other:       param.description,
		})
	}

	return messages
}
