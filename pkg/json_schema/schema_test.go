package json_schema_test

import (
	"encoding/json"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	schema, err := json_schema.New(
		json_schema.Typed(json_schema.TypeString),
		json_schema.Default(json_schema.MustNewValue("test")),
		json_schema.Enum(
			json_schema.MustNewValue("test"),
			json_schema.MustNewValue("test2"),
		),
	)
	require.NoError(t, err)

	schemaBytes, err := json.Marshal(schema)
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"string","default":"test","enum":["test","test2"]}`, string(schemaBytes))

	schemaValue, err := json_schema.NewValue(schema)
	require.NoError(t, err)

	marshaledSchemaBytes, err := schemaValue.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{"type":"string","default":"test","enum":["test","test2"]}`, string(marshaledSchemaBytes))
}
