package config

import (
	"encoding/json"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	"github.com/bitmagnet-io/plugin-qbittorrent/i18n"
)

type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type param struct {
	description string
	schema      json_schema.JSONSchema
}

var params = map[string]param{
	"url": {
		description: "qBittorrent web UI URL",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
			json_schema.Default(json_schema.MustNewValue("http://localhost:8080")),
		),
	},
	"username": {
		description: "qBittorrent username",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
			json_schema.Default(json_schema.MustNewValue("admin")),
		),
	},
	"password": {
		description: "qBittorrent password",
		schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
			json_schema.Default(json_schema.MustNewValue("adminadmin")),
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
