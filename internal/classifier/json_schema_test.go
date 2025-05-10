package classifier

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed json-schema.draft-07.json
var metaSchemaJSON []byte

func TestJSONSchema(t *testing.T) {
	t.Parallel()

	schemaJSON, err := DefaultJSONSchema().MarshalJSON()
	require.NoError(t, err)

	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	metaSchemaLoader := gojsonschema.NewBytesLoader(metaSchemaJSON)

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
}
