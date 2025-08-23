package param_test

import (
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLogLevel(t *testing.T) {
	t.Parallel()

	p, err := param.New(
		param.WithEnumValues(level.LevelValues()...),
		param.WithDefault(level.LevelInfo),
	)

	require.NoError(t, err)
	assert.Equal(t, level.LevelInfo, p.NewDefault())
	assert.Equal(t, "info", p.Stringify(level.LevelInfo))

	parsed, err := p.Parse("info")
	require.NoError(t, err)
	assert.Equal(t, level.LevelInfo, parsed)

	expectEncode := []byte("info\n")
	node, err := p.EncodeYAML(level.LevelInfo)
	require.NoError(t, err)
	encoded, err := yaml.Marshal(node)
	require.NoError(t, err)
	assert.Equal(t, expectEncode, encoded)

	_, err = p.Parse("invalid_level")
	require.Error(t, err)

	// decoded, err := param.DecodeYAML(encoded)
	// require.NoError(t, err)
	// assert.Equal(t, level.LevelInfo, decoded)

	require.NoError(t, p.Validate(level.LevelInfo))
	require.ErrorIs(t, p.Validate("invalid_level"), param.ErrInvalid)

	assert.True(t, p.Equals(level.LevelInfo, level.LevelInfo))
}

func TestUint32(t *testing.T) {
	t.Parallel()

	p, err := param.New(
		param.WithEnumValues(uint32(1), uint32(2), uint32(3)),
		param.WithDefault(uint32(2)),
	)

	require.NoError(t, err)
	assert.Equal(t, uint32(2), p.NewDefault())
	assert.Equal(t, "2", p.Stringify(uint32(2)))

	parsed, err := p.Parse("2")
	require.NoError(t, err)
	assert.Equal(t, uint32(2), parsed)

	expectEncode := []byte("2\n")
	node, err := p.EncodeYAML(uint32(2))
	require.NoError(t, err)
	encoded, err := yaml.Marshal(node)
	require.NoError(t, err)
	assert.Equal(t, expectEncode, encoded)

	decoded, err := p.DecodeYAML(node)
	require.NoError(t, err)
	assert.Equal(t, uint32(2), decoded)

	require.NoError(t, p.Validate(uint32(2)))
	require.ErrorIs(t, p.Validate(uint32(4)), param.ErrInvalid)
}

func TestStringSlice(t *testing.T) {
	t.Parallel()

	p, err := param.New(
		param.WithSlice(
			param.WithStringLiteral[string](),
			param.MaxLength[any, any, string](6),
		),
		param.MinLength[any, string, []string](2),
		param.WithNewDefault(func() []string { return []string{"default1", "default2"} }),
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"default1", "default2"}, p.NewDefault())

	parsed, err := p.Parse("item1,item2")
	require.NoError(t, err)
	assert.Equal(t, []string{"item1", "item2"}, parsed)

	expectEncode := []byte("- item1\n- item2\n")
	node, err := p.EncodeYAML([]string{"item1", "item2"})
	require.NoError(t, err)
	encoded, err := yaml.Marshal(node)
	require.NoError(t, err)
	assert.Equal(t, expectEncode, encoded)

	decoded, err := p.DecodeYAML(node)
	require.NoError(t, err)
	assert.Equal(t, []string{"item1", "item2"}, decoded)

	require.NoError(t, p.Validate([]string{"item1", "item2"}))
	require.ErrorIs(t, p.Validate([]string{"item1"}), param.ErrInvalid)
	require.ErrorIs(t, p.Validate([]string{"item1", "item2verylongstring"}), param.ErrInvalid)
}

func TestStruct(t *testing.T) {
	type TestStruct struct {
		CamelCase string
		Bar       int
	}

	p, err := param.New(
		param.WithMapstructure[TestStruct](),
		param.WithNewDefault(func() TestStruct {
			return TestStruct{
				CamelCase: "test",
				Bar:       2,
			}
		}),
	)

	require.NoError(t, err)

	assert.Equal(t, TestStruct{
		CamelCase: "test",
		Bar:       2,
	}, p.NewDefault())

	yamlBytes := []byte(`{"camel_case":"foo","bar":6}`)
	var yamlNode yaml.Node
	require.NoError(t, yaml.Unmarshal(yamlBytes, &yamlNode))

	decoded, err := p.DecodeYAML(yamlNode)
	require.NoError(t, err)
	assert.Equal(t, TestStruct{
		CamelCase: "foo",
		Bar:       6,
	}, decoded)
}
