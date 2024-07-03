package servarr

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
)

type Prowlarr struct {
	download *servarrDownload
}

func (prowlarr Prowlarr) SearchQueryParam() (map[string]string, error) {
	qp := map[string]string{
		"query":      prowlarr.download.content.InfoHash.String(),
		"indexerIds": fmt.Sprint(prowlarr.download.indexer.ID),
		"type":       "search",
	}

	return qp, nil
}

func NewProwlarr(content *gqlmodel.TorrentContent, config *Config, ctx context.Context) *servarrDownload {
	if config.Prowlarr.ApiKey != "private" {
		dl := &servarrDownload{
			ctx:            ctx,
			config:         config,
			content:        content,
			api:            config.Prowlarr,
			searchEndpoint: "search",
			apiVersion:     "v1",
			indexerName:    config.IndexerName,
		}
		dl.arr = Prowlarr{download: dl}
		return dl
	}
	return nil
}
