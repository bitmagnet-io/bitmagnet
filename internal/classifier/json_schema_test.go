package classifier

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestJSONSchema(t *testing.T) {
	t.Parallel()

	defaultSchema := DefaultJSONSchema()

	schemaBytes, err := json.MarshalIndent(defaultSchema, "", "  ")
	require.NoError(t, err)

	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	metaSchemaLoader := gojsonschema.NewBytesLoader(json_spec.MetaSchema)

	// validate the schema against the meta schema
	metaResult, err := gojsonschema.Validate(metaSchemaLoader, schemaLoader)
	require.NoError(t, err)
	assert.True(t, metaResult.Valid())

	coreClassifier, err := yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}.source()
	require.NoError(t, err)
	coreClassifierJSON, err := json.Marshal(coreClassifier)
	require.NoError(t, err)

	documentLoader := gojsonschema.NewBytesLoader(coreClassifierJSON)

	// validate the classifier against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	require.NoError(t, err)
	assert.True(t, result.Valid())

	var unmarshaledSchema json_schema.JSONSchema
	require.NoError(t, json.Unmarshal(schemaBytes, &unmarshaledSchema))
}
