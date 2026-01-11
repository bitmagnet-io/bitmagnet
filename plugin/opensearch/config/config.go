//go:build wasip1

package config

type Config struct {
	Addresses   []string `json:"addresses"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	IndexPrefix string   `json:"index_prefix"`
}
