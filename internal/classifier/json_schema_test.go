package classifier

import (
	_ "embed"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"testing"
)

//go:embed json-schema.draft-07.json
var metaSchemaJson []byte

func TestJsonSchema(t *testing.T) {

	schemaJson, err := DefaultJsonSchema().MarshalJSON()
	assert.NoError(t, err)

	schemaLoader := gojsonschema.NewBytesLoader(schemaJson)
	metaSchemaLoader := gojsonschema.NewBytesLoader(metaSchemaJson)

	// validate the schema against the meta schema
	metaResult, err := gojsonschema.Validate(metaSchemaLoader, schemaLoader)
	assert.NoError(t, err)
	assert.True(t, metaResult.Valid())

	coreClassifier, err := yamlSourceProvider{rawSourceProvider: coreSourceProvider{}}.source()
	assert.NoError(t, err)
	coreClassifierJson, err := json.Marshal(coreClassifier)

	documentLoader := gojsonschema.NewBytesLoader(coreClassifierJson)

	// validate the classifier against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	assert.NoError(t, err)
	assert.True(t, result.Valid())
}
