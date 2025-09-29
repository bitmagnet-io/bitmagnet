package json_schema_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	t.Parallel()

	yamlBytes := []byte(`foo: "bar"` + "\n")

	var value json_schema.JSONValue

	require.NoError(t, value.UnmarshalYAML(yamlBytes))

	jsonBytes, err := value.MarshalJSON()
	require.NoError(t, err)

	assert.Equal(t, []byte(`{"foo":"bar"}`), jsonBytes)

	mYamlBytes, err := value.MarshalYAML()
	require.NoError(t, err)

	assert.Equal(t, yamlBytes, mYamlBytes)
}
