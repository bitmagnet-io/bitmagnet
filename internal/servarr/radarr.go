package servarr

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Radarr struct {
	download *servarrDownload
}

func (radarr Radarr) SearchQueryParam() (map[string]string, error) {
	var movie []MovieResource
	resp, err := radarr.download.resty(&movie).SetQueryParam("tmdbId", radarr.download.content.Content.ID).Get(radarr.download.url("movie"))
	if err != nil {
		return nil, err
	} else if len(movie) < 1 {
		return nil, radarr.download.fmtErr("movie not found", resp)
	}

	qp := map[string]string{
		"movieId": fmt.Sprint(movie[0].ID),
	}

	return qp, nil
}

func NewRadarr(content *gqlmodel.TorrentContent, config *Config, ctx context.Context) *servarrDownload {
	if content.ContentType.ContentType == model.ContentTypeMovie &&
		config.Radarr.ApiKey != "private" {
		dl := &servarrDownload{
			ctx:            ctx,
			config:         config,
			content:        content,
			api:            config.Radarr,
			searchEndpoint: "release",
			apiVersion:     "v3",
			indexerName:    fmt.Sprintf("%s (Prowlarr)", config.IndexerName),
		}
		dl.arr = Radarr{download: dl}
		return dl
	}
	return nil
}
