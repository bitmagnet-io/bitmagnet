package postgres

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"go.uber.org/fx"
)

type (
	DSN      string
	Host     string
	Port     int
	User     string
	Password string
	Database string
)

var (
	ParamDSN = param.MustNew[DSN]()

	ParamHost = param.MustNew(
		param.WithDefault(Host("localhost")),
	)

	ParamPort = param.MustNew(
		param.WithDefault(Port(5432)),
	)

	ParamUser = param.MustNew(
		param.WithDefault(User("postgres")),
	)

	ParamPassword = param.MustNew(
		param.WithDefault(Password("postgres")),
	)

	ParamDatabase = param.MustNew(
		param.WithDefault(Database("bitmagnet")),
	)
)

type Config struct {
	fx.In
	DSN      DSN
	Host     Host
	Port     Port
	User     User
	Password Password
	Database Database
}

func (c Config) CreateDSN() DSN {
	if c.DSN != "" {
		return c.DSN
	}

	var values []string

	for k, v := range c.values() {
		values = append(values, fmt.Sprintf("%s=%s", k, v))
	}

	slices.Sort(values)

	return DSN(strings.Join(values, " "))
}

func (c Config) values() map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", c.Database)
	setIfNotEmpty(p, "user", c.User)
	setIfNotEmpty(p, "host", c.Host)
	setIfNotEmpty(p, "port", fmt.Sprintf("%d", c.Port))
	setIfNotEmpty(p, "password", c.Password)

	return p
}

func setIfNotEmpty(m map[string]string, key string, val any) {
	strVal := fmt.Sprintf("%v", val)
	if strVal != "" {
		m[key] = strVal
	}
}
