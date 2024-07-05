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
	"go.uber.org/zap"
)

type servarrSpecific interface {
	SearchQueryParam() (map[string]string, error)
}

type servarrDownload struct {
	ctx                 context.Context
	config              *Config
	logger              *zap.SugaredLogger
	content             *gqlmodel.TorrentContent
	api                 UrlKey
	arr                 servarrSpecific
	searchEndpoint      string
	apiVersion          string
	indexerName         string
	onlySearchBitmagnet bool
	indexers            []IndexerResource
	indexer             *IndexerResource
}

type ServarrClient struct {
	Config *Config
	Search *search.Search
	Logger *zap.SugaredLogger
}

func (d *servarrDownload) resty(res interface{}) *resty.Request {
	return resty.
		New().
		R().
		SetContext(d.ctx).
		SetLogger(d.logger).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Api-Key", d.api.ApiKey).
		SetResult(res)
}

func (d *servarrDownload) url(path string) string {
	return fmt.Sprintf("%s/api/%s/%s", d.api.Url, d.apiVersion, path)
}

func (d *servarrDownload) fmtErr(msg string, resp *resty.Response, err error) error {
	if resp != nil {
		msg = fmt.Sprintf("%s [%s] (%d)", msg, resp.Request.URL, resp.StatusCode())
	}
	err = errors.Join(errors.New(msg), err)
	d.logger.Warn(err.Error())
	return err
}

func (d *servarrDownload) setIndexerEnabled(enable bool) error {
	var ids []*int64

	for _, i := range d.indexers {
		if i.EnableInteractiveSearch && i.Name != d.indexerName {
			ids = append(ids, &i.ID)
		}
	}

	if len(ids) > 0 {
		resp, err := d.resty(make([]IndexerResource, 0)).
			SetBody(&IndexerBulkResource{Ids: ids, EnableInteractiveSearch: &enable}).
			Put(d.url("indexer/bulk"))
		if err != nil || resp.IsError() {
			return d.fmtErr("indexer bulk", resp, err)
		}
	}

	return nil

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
	if err != nil || resp.IsError() {
		return d.fmtErr("get indexer", resp, err)
	}

	i := slices.IndexFunc(d.indexers, func(i IndexerResource) bool { return i.Name == d.indexerName })
	if i == -1 {
		return d.fmtErr(fmt.Sprintf("indexer not found %s", d.config.IndexerName), resp, nil)
	}
	d.indexer = &d.indexers[i]

	qp, err := transformSearchQueryParam(d.arr)
	if err != nil {
		return err
	}

	if d.onlySearchBitmagnet {
		// any indexers that are disabled, defer re-enablement
		defer d.setIndexerEnabled(true)
		err = d.setIndexerEnabled(false)
		if err != nil {
			return err
		}
	}

	resp, err = d.resty(&releases).SetQueryParams(qp).Get(d.url(d.searchEndpoint))
	if err != nil || resp.IsError() {
		return d.fmtErr("get download", resp, err)
	}

	i = slices.IndexFunc(releases, func(r ReleaseResource) bool { return r.GUID == d.content.InfoHash.String() })
	if i == -1 {
		return d.fmtErr("download not found", resp, nil)
	}
	resp, err = d.resty(&release).
		SetBody(&ReleaseResource{GUID: releases[i].GUID, IndexerID: d.indexer.ID}).
		Post(d.url(d.searchEndpoint))
	if err != nil || resp.IsError() {
		return d.fmtErr("trigger download", resp, err)
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
		msg := fmt.Sprintf("no servarr client accepts (%s) (%s)", content.ContentType.ContentType, content.InfoHash)
		c.Logger.Warn(msg)
		return errors.New(msg)
	}

	d.logger = c.Logger
	c.Logger.Debugf("%s %T", content.InfoHash.String(), d.arr)

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
