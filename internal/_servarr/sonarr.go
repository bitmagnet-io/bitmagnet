package servarr

import (
	"context"
	"fmt"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Sonarr struct {
	download *servarrDownload
}

func (sonarr Sonarr) SearchQueryParam() (map[string]string, error) {
	i := slices.IndexFunc(sonarr.download.content.Content.Attributes, func(a model.ContentAttribute) bool { return a.Source == "tvdb" })
	if i == -1 {
		return nil, fmt.Errorf("TVDB ID not available")
	}

	var series []SeriesResource
	resp, err := sonarr.download.resty(&series).SetQueryParam("tvdbId",
		sonarr.download.content.Content.Attributes[i].Value).Get(sonarr.download.url("series"))
	if err != nil || resp.IsError() {
		return nil, sonarr.download.fmtErr("get series", resp, err)
	} else if len(series) != 1 {
		return nil, sonarr.download.fmtErr("series not found", resp, nil)
	}

	qp := map[string]string{
		"seriesId":     fmt.Sprint(series[0].ID),
		"seasonNumber": fmt.Sprint(sonarr.download.content.Episodes.Seasons[0].Season),
	}
	var episodes []EpisodeResource
	resp, err = sonarr.download.resty(&episodes).SetQueryParams(qp).Get(sonarr.download.url("episode"))
	if err != nil || resp.IsError() {
		return nil, sonarr.download.fmtErr("get episode", resp, err)
	}
	i = slices.IndexFunc(episodes, func(e EpisodeResource) bool {
		return e.SeasonNumber == int64(sonarr.download.content.Episodes.Seasons[0].Season) &&
			e.EpisodeNumber == int64(sonarr.download.content.Episodes.Seasons[0].Episodes[0])
	})
	if i == -1 {
		return nil, sonarr.download.fmtErr("episode not found", resp, nil)
	}

	qp["episodeId"] = fmt.Sprint(episodes[i].ID)

	return qp, nil
}

func NewSonarr(content *gqlmodel.TorrentContent, config *Config, ctx context.Context) *servarrDownload {
	if content.ContentType.ContentType == model.ContentTypeTvShow &&
		len(content.Episodes.Seasons) == 1 &&
		len(content.Episodes.Seasons[0].Episodes) == 1 &&
		config.Sonarr.ApiKey != "private" {
		dl := &servarrDownload{
			ctx:                 ctx,
			config:              config,
			content:             content,
			api:                 config.Sonarr,
			searchEndpoint:      "release",
			apiVersion:          "v3",
			onlySearchBitmagnet: config.OnlySearchBitmagnet,
			indexerName:         fmt.Sprintf("%s (Prowlarr)", config.IndexerName),
		}
		dl.arr = Sonarr{download: dl}
		return dl
	}
	return nil
}
