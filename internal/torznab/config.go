package torznab

import (
	"reflect"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
)

type Profile struct {
	Name           string
	OrderBy        search.TorrentContentOrderBy
	OrderDirection search.OrderDirection
	Tags           []string
}

type Config struct {
	// Profiles       []Profile          // config framework does not support stuct,slice,struct,item+
	// Profiles       map[string]Profile // config framework does not support struct,map[string],struct,item+
	DefaultProfile Profile
	Profile0       Profile
	Profile1       Profile
	Profile2       Profile
	Profile3       Profile
	Profile4       Profile
	Hostname       string
}

func defaultProfile(name string) Profile {
	return Profile{
		Name:           name,
		OrderBy:        search.TorrentContentOrderByRelevance,
		OrderDirection: search.OrderDirectionDescending,
	}
}

func NewDefaultConfig() Config {
	config := Config{}
	// create a default profile for each Profile defined in Config struct
	configValue := reflect.ValueOf(config)
	for i := 0; i < configValue.Type().NumField(); i++ {
		field := configValue.Type().Field(i)
		if field.Type == reflect.ValueOf(Profile{}).Type() {
			reflect.ValueOf(&config).Elem().FieldByName(field.Name).Set(
				reflect.ValueOf(defaultProfile(strings.ToLower(field.Name))),
			)
		}
	}
	return config

}

func (configPtr *Config) Map() map[string]Profile {
	configValue := reflect.ValueOf(*configPtr)
	// this will create an allocation larger than required
	configMap := make(map[string]Profile, configValue.Type().NumField())
	for i := 0; i < configValue.Type().NumField(); i++ {
		field := configValue.Type().Field(i)
		if field.Type == reflect.ValueOf(Profile{}).Type() {
			profile := reflect.ValueOf(configPtr).Elem().Field(i).Interface().(Profile)
			configMap[profile.Name] = profile
		}
	}

	return configMap
}
