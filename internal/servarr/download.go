package servarr

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/client"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/go-resty/resty/v2"
)

type servarrSpecific interface {
	SearchQueryParam() (map[string]string, error)
}

type servarrDownload struct {
	ctx            context.Context
	config         *Config
	content        *gqlmodel.TorrentContent
	api            UrlKey
	arr            servarrSpecific
	searchEndpoint string
	apiVersion     string
	indexerName    string
	indexers       []IndexerResource
	indexer        *IndexerResource
}

type ServarrClient struct {
	Config *Config
	Search *search.Search
}

func (d *servarrDownload) resty(res interface{}) *resty.Request {
	return resty.New().R().SetContext(d.ctx).SetHeader("Content-Type", "application/json").SetHeader("X-Api-Key", d.api.ApiKey).SetResult(res)
}

func (d *servarrDownload) url(path string) string {
	return fmt.Sprintf("%s/api/%s/%s", d.api.Url, d.apiVersion, path)
}

func (d *servarrDownload) fmtErr(msg string, resp *resty.Response) error {
	return fmt.Errorf("%s [%s] (%d)", msg, resp.Request.URL, resp.StatusCode())
}

func transformSearchQueryParam(arr servarrSpecific) (map[string]string, error) {
	return arr.SearchQueryParam()
}

func (d *servarrDownload) download() error {
	var (
		releases []ReleaseResource
		release  ReleaseResource
	)
	resp, err := d.resty(&d.indexers).Get(d.url("indexer"))
	if err != nil {
		return err
	} else if resp.IsError() {
		return d.fmtErr("http request failed", resp)
	}

	i := slices.IndexFunc(d.indexers, func(i IndexerResource) bool { return i.Name == d.indexerName })
	if i == -1 {
		return d.fmtErr(fmt.Sprintf("indexer not found %s", d.config.IndexerName), resp)
	}
	d.indexer = &d.indexers[i]

	qp, err := transformSearchQueryParam(d.arr)
	if err != nil {
		return err
	}
	resp, err = d.resty(&releases).SetQueryParams(qp).Get(d.url(d.searchEndpoint))
	if err != nil {
		return err
	}

	i = slices.IndexFunc(releases, func(r ReleaseResource) bool { return r.GUID == d.content.InfoHash.String() })
	if i == -1 {
		return d.fmtErr("download not found", resp)
	}
	_, err = d.resty(&release).SetBody(&ReleaseResource{GUID: releases[i].GUID, IndexerID: d.indexer.ID}).Post(d.url(d.searchEndpoint))
	if err != nil {
		return err
	}

	return nil
}

func (c *ServarrClient) downloadOne(ctx context.Context, content *gqlmodel.TorrentContent) error {
	d := NewSonarr(content, c.Config, ctx)
	if d == nil {
		d = NewRadarr(content, c.Config, ctx)
	}
	if d == nil {
		d = NewProwlarr(content, c.Config, ctx)
	}
	if d == nil {
		return fmt.Errorf("no servarr client accepts (%s) (%s)", content.ContentType.ContentType, content.InfoHash)
	}

	return d.download()

}

func (c *ServarrClient) AddInfoHashes(ctx context.Context, req client.AddInfoHashesRequest) error {
	var allErr error
	for _, infoHash := range req.InfoHashes {
		content, err := gqlmodel.TorrentContentQuery{TorrentContentSearch: *c.Search}.Search(
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
		err = c.downloadOne(ctx, &content.Items[0])
		allErr = errors.Join(allErr, err)

	}
	return allErr
}
