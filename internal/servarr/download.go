package servarr

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/go-resty/resty/v2"
)

type downloadHelper struct {
	ctx    context.Context
	config *Config
	tc     *gqlmodel.TorrentContent
	api    UrlKey
	use    string
}

func (d *downloadHelper) resty(res interface{}) *resty.Request {
	return resty.New().R().SetContext(d.ctx).SetHeader("Content-Type", "application/json").SetHeader("X-Api-Key", d.api.ApiKey).SetResult(res)
}

func (d *downloadHelper) url(path string) string {
	v := "v3"
	if d.use == "prowlarr" {
		v = "v1"
	}
	return fmt.Sprintf("%s/api/%s/%s", d.api.Url, v, path)
}

func (d *downloadHelper) sonarrreleases() ([]ReleaseResource, error) {
	var releases []ReleaseResource
	i := slices.IndexFunc(d.tc.Content.Attributes, func(a model.ContentAttribute) bool { return a.Source == "tvdb" })
	if i == -1 {
		return nil, fmt.Errorf("TVDB ID not available")
	}

	var series []SeriesResource
	_, err := d.resty(&series).SetQueryParam("tvdbId", d.tc.Content.Attributes[i].Value).Get(d.url("series"))
	if err != nil {
		return nil, err
	} else if len(series) != 1 {
		return nil, fmt.Errorf("series not found in sonarr %s", d.tc.Content.Attributes[i].Value)
	}

	qp := map[string]string{
		"seriesId":     fmt.Sprint(series[0].ID),
		"seasonNumber": fmt.Sprint(d.tc.Episodes.Seasons[0].Season),
	}
	var episodes []EpisodeResource
	_, err = d.resty(&episodes).SetQueryParams(qp).Get(d.url("episode"))
	if err != nil {
		return nil, err
	}
	i = slices.IndexFunc(episodes, func(e EpisodeResource) bool {
		return e.SeasonNumber == int64(d.tc.Episodes.Seasons[0].Season) && e.EpisodeNumber == int64(d.tc.Episodes.Seasons[0].Episodes[0])
	})
	if i == -1 {
		return nil, fmt.Errorf("Episode not found in sonarr %s", d.tc.Content.Attributes[i].Value)
	}

	qp["episodeId"] = fmt.Sprint(episodes[i].ID)
	_, err = d.resty(&releases).SetQueryParams(qp).Get(d.url("release"))
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (d *downloadHelper) radarrreleases() ([]ReleaseResource, error) {
	var releases []ReleaseResource

	var movie []MovieResource
	_, err := d.resty(&movie).SetQueryParam("tmdbId", d.tc.Content.ID).Get(d.url("movie"))
	if err != nil {
		return nil, err
	} else if len(movie) < 1 {
		return nil, fmt.Errorf("movie not found in radarr (%s)", d.tc.Content.ID)
	}

	_, err = d.resty(&releases).SetQueryParam("movieId", fmt.Sprint(movie[0].ID)).Get(d.url("release"))
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (d *downloadHelper) download() (*string, error) {
	var (
		indexers []IndexerResource
		releases []ReleaseResource
		release  ReleaseResource
	)
	endpoint := "release"
	resp, err := d.resty(&indexers).Get(d.url("indexer"))
	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("%s http request failed (%d) [%s]", d.use, resp.StatusCode(), resp.Status())
	}

	name := fmt.Sprintf("%s (Prowlarr)", d.config.IndexerName)
	if d.use == "prowlarr" {
		name = d.config.IndexerName
	}

	i := slices.IndexFunc(indexers, func(i IndexerResource) bool { return i.Name == name })
	if i == -1 {
		return nil, fmt.Errorf("indexer not found %s first %s", d.config.IndexerName, indexers[0].Name)
	}
	indexer := indexers[i]

	switch d.use {
	case "sonarr":
		releases, err = d.sonarrreleases()
		if err != nil {
			return nil, err
		}
	case "radarr":
		releases, err = d.radarrreleases()
		if err != nil {
			return nil, err
		}

	default:
		endpoint = "search"
		qp := map[string]string{
			"query":      d.tc.InfoHash.String(),
			"indexerIds": fmt.Sprint(indexer.ID),
			"type":       "search",
		}
		_, err = d.resty(&releases).SetQueryParams(qp).Get(d.url(endpoint))
		if err != nil {
			return nil, err
		}

	}

	i = slices.IndexFunc(releases, func(r ReleaseResource) bool { return r.GUID == d.tc.InfoHash.String() })
	if i == -1 {
		return nil, fmt.Errorf("download not found")
	}
	_, err = d.resty(&release).SetBody(&ReleaseResource{GUID: releases[i].GUID, IndexerID: indexer.ID}).Post(d.url(endpoint))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func downloadOne(ctx context.Context, tc *gqlmodel.TorrentContent, config *Config) (*string, error) {
	dltype := "prowlarr"
	api := config.Prowlarr
	if tc.ContentType.ContentType.String() == "tv_show" && len(tc.Episodes.Seasons) == 1 && len(tc.Episodes.Seasons[0].Episodes) == 1 && config.Sonarr.ApiKey != "private" {
		dltype = "sonarr"
		api = config.Sonarr
	} else if tc.ContentType.ContentType.String() == "movie" && config.Radarr.ApiKey != "private" {
		dltype = "radarr"
		api = config.Radarr
	}

	if api.ApiKey == "private" {
		return nil, fmt.Errorf("no ApiKey for %s", dltype)
	}

	d := &downloadHelper{
		ctx:    ctx,
		config: config,
		tc:     tc,
		api:    api,
		use:    dltype,
	}

	return d.download()

}

func Download(ctx context.Context, config *Config, search search.Search, infoHashes []protocol.ID) (*string, error) {
	var allErr error
	for _, infoHash := range infoHashes {
		content, err := gqlmodel.TorrentContentQuery{TorrentContentSearch: search}.Search(
			ctx,
			&query.SearchParams{QueryString: model.NullString{String: infoHash.String(), Valid: true}},
			&gen.TorrentContentFacetsInput{},
			make([]gen.TorrentContentOrderByInput, 0),
		)
		if err != nil {
			allErr = errors.Join(allErr, err)
		} else if len(content.Items) != 1 {
			allErr = errors.Join(allErr, fmt.Errorf("Too many content results (%d) for download", len(content.Items)))
		}
		_, err = downloadOne(ctx, &content.Items[0], config)
		allErr = errors.Join(allErr, err)

	}
	return nil, allErr
}
