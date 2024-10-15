package postgres

import (
	"fmt"
	"strings"
)

type Config struct {
	DSN               string
	Host              string
	User              string
	Port              uint
	Name              string
	ConnectionTimeout uint
	Password          string
	SSLMode           string
	SSLCertPath       string
	SSLKeyPath        string
	SSLRootCertPath   string
}

func NewDefaultConfig() Config {
	return Config{
		Host: "localhost",
		User: "postgres",
		Port: 5432,
		Name: "bitmagnet",
	}
}

func (c *Config) CreateDSN() string {
	if c.DSN != "" {
		return c.DSN
	}
	vals := dbValues(c)
	p := make([]string, 0, len(vals))
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key string, val interface{}) {
	strVal := fmt.Sprintf("%v", val)
	if strVal != "" {
		m[key] = strVal
	}
}

func setIfPositive(m map[string]string, key string, val uint) {
	if val > 0 {
		m[key] = fmt.Sprintf("%d", val)
	}
}

func dbValues(cfg *Config) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", cfg.Name)
	setIfNotEmpty(p, "user", cfg.User)
	setIfNotEmpty(p, "host", cfg.Host)
	setIfNotEmpty(p, "port", fmt.Sprintf("%d", cfg.Port))
	setIfNotEmpty(p, "sslmode", cfg.SSLMode)
	setIfPositive(p, "connect_timeout", cfg.ConnectionTimeout)
	setIfNotEmpty(p, "password", cfg.Password)
	setIfNotEmpty(p, "sslcert", cfg.SSLCertPath)
	setIfNotEmpty(p, "sslkey", cfg.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", cfg.SSLRootCertPath)
	return p
}
