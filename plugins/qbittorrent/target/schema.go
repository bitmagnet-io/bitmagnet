//go:build wasip1

package target

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_forms"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	"github.com/bitmagnet-io/plugin-qbittorrent/i18n"
)

type Data struct {
	Category string `json:"category,omitempty"`
	Stopped  bool   `json:"stopped"`
}

func (t *target) DataSchema(ctx context.Context, empty *api.Empty) (*api.JSONPayload, error) {
	categories, err := t.getCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	schemaBytes, err := json.Marshal(newDataSchema(categories))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data schema: %w", err)
	}

	return &api.JSONPayload{Data: schemaBytes}, nil
}

func (t *target) getCategories(ctx context.Context) ([]string, error) {
	var categoriesMap map[string]any

	err := t.retryWithLogin(ctx, func() error {
		return t.request(ctx, &http.Request{
			Url:    t.config.URL + "/api/v2/torrents/categories",
			Method: http.Method_get,
		}, &categoriesMap)
	})
	if err != nil {
		return nil, err
	}

	return slices.Sorted(maps.Keys(categoriesMap)), nil
}

func (t *target) UISchema(ctx context.Context, params *api.SendTorrentsUISchemaParams) (*api.JSONPayload, error) {
	schemaBytes, err := json.Marshal(newUISchema(params.AcceptLanguage))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UI schema: %w", err)
	}

	return &api.JSONPayload{Data: schemaBytes}, nil
}

func newDataSchema(categories []string) json_schema.JSONSchema {
	properties := map[string]json_schema.JSONSchema{
		"stopped": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeBoolean),
			json_schema.Default(json_schema.MustNewValue(false)),
		),
	}

	if len(categories) > 0 {
		nullableCategories := make([]json_schema.JSONValue, len(categories)+1)
		nullableCategories[0] = json_schema.MustNewValue(nil)
		for i, category := range categories {
			nullableCategories[i+1] = json_schema.MustNewValue(category)
		}
		properties["category"] = json_schema.MustNew(
			json_schema.Enum(nullableCategories...),
		)
	}

	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeObject),
		json_schema.Properties(properties),
		json_schema.AdditionalPropertiesFalse(),
		json_schema.Default(json_schema.MustNewValue(map[string]any{})),
	)
}

func newUISchema(acceptLanguage []string) json_forms.UISchema {
	localizer := i18n.NewLocalizer(acceptLanguage)

	return json_forms.Layout{
		Type: json_forms.LayoutTypeVertical,
		Elements: []json_forms.Element{
			json_forms.Control{
				Type:  "Control",
				Scope: "#/properties/category",
				Label: ptr(localizer.Localize("category")),
			},
			json_forms.Control{
				Type:  "Control",
				Scope: "#/properties/stopped",
				Label: ptr(localizer.Localize("stopped")),
			},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
