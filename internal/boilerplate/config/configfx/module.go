package configfx

import (
	"github.com/adrg/xdg"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configresolver"
	"go.uber.org/fx"
	"os"
	"strings"
)

func New() fx.Option {
	osEnv := ReadOsEnv()

	var options []fx.Option

	var extraConfigFiles []string
	if osEnv[extraFilesKey] != "" {
		extraConfigFiles = strings.Split(osEnv[extraFilesKey], ",")
	}
	for i, file := range extraConfigFiles {
		options = append(options,
			fx.Provide(
				fx.Annotated{
					Group: "config_resolvers",
					Target: func() (configresolver.Resolver, error) {
						return configresolver.NewFromYamlFile(
							file,
							false,
							configresolver.WithPriority(-i),
						)
					},
				},
			))
	}

	options = append(options,
		fx.Provide(config.New),
		fx.Provide(fx.Annotated{
			Group: "config_resolvers",
			Target: func() (configresolver.Resolver, error) {
				return configresolver.NewEnv(
					osEnv,
					configresolver.WithPriority(-len(extraConfigFiles)),
				), nil
			},
		}),
		fx.Provide(
			fx.Annotated{
				Group: "config_resolvers",
				Target: func() (configresolver.Resolver, error) {
					return configresolver.NewFromYamlFile(
						"./config.yml",
						true,
						configresolver.WithPriority(10),
					)
				},
			},
		),
	)
	if configFilePath, err := xdg.ConfigFile("bitmagnet/config.yml"); err == nil {
		options = append(options,
			fx.Provide(
				fx.Annotated{
					Group: "config_resolvers",
					Target: func() (configresolver.Resolver, error) {
						return configresolver.NewFromYamlFile(
							configFilePath,
							true,
							configresolver.WithPriority(20),
						)
					},
				},
			),
		)
	}
	return fx.Module(
		"config",
		fx.Options(options...),
	)
}

func ReadOsEnv() map[string]string {
	rawEnv := os.Environ()
	env := make(map[string]string, len(rawEnv))
	for _, rawEnvEntry := range rawEnv {
		parts := strings.SplitN(rawEnvEntry, "=", 2)
		env[parts[0]] = parts[1]
	}
	return env
}

const extraFilesKey = "EXTRA_CONFIG_FILES"
