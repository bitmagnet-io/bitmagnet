package client

import (
	"github.com/bitmagnet-io/bitmagnet/proto/host/http_client"
	"github.com/bitmagnet-io/plugin-opensearch/config"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func New(httpClient http_client.Service, cfg config.Config) (*opensearchapi.Client, error) {
	return opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Transport: &transport{
				httpClient: httpClient,
			},
			Addresses: cfg.Addresses,
			Username:  cfg.Username,
			Password:  cfg.Password,
		},
	})
}
