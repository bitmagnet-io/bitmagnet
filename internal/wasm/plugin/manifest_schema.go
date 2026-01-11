package plugin

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/xeipuuv/gojsonschema"
)

const schemaID = "https://bitmagnet.io/schemas/plugin-manifest-0.1.json"

var ManifestSchema = json_schema.MustNew(
	json_schema.MetaschemaDraft7(),
	json_schema.Typed(json_schema.TypeObject),
	json_schema.Properties(map[string]json_schema.JSONSchema{
		"$schema": json_schema.MustNew(
			json_schema.Const(json_schema.MustNewValue(schemaID)),
		),
		"name": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
		),
		"description": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
		),
		"version": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
		),
		"concurrency": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeInteger),
			json_schema.Minimum(1),
			json_schema.Maximum(100),
			json_schema.Default(json_schema.MustNewValue(1)),
		),
		"capabilities": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeObject),
			json_schema.Properties(map[string]json_schema.JSONSchema{
				capabilitySearchAdapterName: CapabilitySearchAdapter{}.jsonSchema(),
				capabilityIndexerName:       CapabilityIndexer{}.jsonSchema(),
				capabilityHTTPHandlerName:   CapabilityHTTPHandler{}.jsonSchema(),
				capabilityReceiverName:      CapabilityReceiver{}.jsonSchema(),
			}),
			json_schema.AdditionalPropertiesFalse(),
		),
		"permissions": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeObject),
			json_schema.Properties(map[string]json_schema.JSONSchema{
				"fs":   PermissionFS{}.jsonSchema(),
				"http": PermissionHTTP{}.jsonSchema(),
			}),
			json_schema.AdditionalPropertiesFalse(),
		),
		"config": json_schema.MustNew(
			json_schema.Typed(json_schema.TypeObject),
		),
	}),
	json_schema.Required(json_schema.RequiredFields{
		"name",
		"version",
	}),
	json_schema.AdditionalPropertiesFalse(),
)

var manifestGoJSONSchema = func() *gojsonschema.Schema {
	compiled, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(ManifestSchema))
	if err != nil {
		panic(err)
	}

	return compiled
}()
