package configfx

import (
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/configresolver"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestConfig(t *testing.T) {
	type Nested struct {
		NestedKey           string `validate:"uppercase"`
		NestedKeyFromConfig string
		IntWithValidation   int `validate:"min=1,max=10"`
		IntSlice            []int
		DB                  int
	}

	type TestConfig struct {
		Foo         string `validate:"lowercase"`
		Bar         int
		Nested      Nested
		StringSlice []string
		StructSlice []Nested
		Duration    time.Duration
		Duration2   time.Duration
	}

	defaultConfig := TestConfig{
		Foo:         "default",
		Bar:         2,
		StringSlice: []string{"a", "b", "c"},
		Nested: Nested{
			NestedKey: "NESTED",
		},
	}
	fx.New(
		NewConfigModule[TestConfig]("test", defaultConfig),
		fx.Provide(
			fx.Annotated{
				Group: "config_resolvers",
				Target: func() (configresolver.Resolver, error) {
					return configresolver.NewEnv(map[string]string{
						"TEST_DURATION":                   "2s",
						"TEST_NESTED_INT_WITH_VALIDATION": "2",
						"TEST_NESTED_DB":                  "3",
					}, configresolver.WithPriority(-10)), nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Group: "config_resolvers",
				Target: func(val *validator.Validate) (configresolver.Resolver, error) {
					return configresolver.NewFromYamlFile("./test_config.yaml", false, val)
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Group: "config_resolvers",
				Target: func(val *validator.Validate) (configresolver.Resolver, error) {
					return configresolver.NewFromYamlFile("./missing.yaml", true, val)
				},
			},
		),
		fx.Provide(validator.New),
		fx.Provide(config.New),
		fx.Invoke(func(cfg TestConfig, shutdowner fx.Shutdowner) {
			assert.Equal(t, TestConfig{
				Foo:         "foo",
				Bar:         2,
				StringSlice: []string{"a", "b", "c"},
				Duration:    time.Second * 2,
				Duration2:   time.Second * 3,
				Nested: Nested{
					NestedKey:           "NESTED",
					NestedKeyFromConfig: "from_config",
					IntWithValidation:   2,
					IntSlice:            []int{1, 2, 3},
					DB:                  3,
				},
				StructSlice: []Nested{
					{
						NestedKey:         "FOO",
						IntWithValidation: 1,
					},
					{
						NestedKey:         "BAR",
						IntWithValidation: 2,
					},
				},
			}, cfg)
			shutdownErr := shutdowner.Shutdown()
			assert.NoError(t, shutdownErr)
		}),
	).Run()
}
