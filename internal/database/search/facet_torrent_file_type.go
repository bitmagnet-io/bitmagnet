package search

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
	"strings"
	"sync"
)

const TorrentFileTypeFacetKey = "file_type"

func TorrentFileTypeFacet(options ...query.FacetOption) query.Facet {
	return torrentFileTypeFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentFileTypeFacetKey),
				query.FacetHasLabel("File Type"),
				query.FacetUsesOrLogic(),
				query.FacetHasAggregationOption(query.RequireJoin(model.TableNameTorrentContent)),
			}, options...)...,
		),
	}
}

type torrentFileTypeFacet struct {
	query.FacetConfig
}

func (f torrentFileTypeFacet) Aggregate(ctx query.FacetContext) (query.AggregationItems, error) {
	type result struct {
		FileType model.FileType
		Count    uint
	}
	var allExts []string
	ftFieldTemplate := "case "
	for _, ft := range model.FileTypeValues() {
		exts := ft.Extensions()
		allExts = append(allExts, exts...)
		ftFieldTemplate += "when {table}.extension in " + makeStringList(exts...) + " then '" + ft.String() + "' "
	}
	ftFieldTemplate += "end as file_type"
	var fileResults []result
	var torrentResults []result
	var errs []error
	// we need to gather aggregations from both the torrent_files table and the torrents table (for the case when the torrent is a single file)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		q, qErr := ctx.NewAggregationQuery(
			query.Table(model.TableNameTorrentFile),
			query.Join(func(q *dao.Query) []query.TableJoin {
				return []query.TableJoin{
					{
						Table: q.TorrentContent,
						On: []field.Expr{
							q.TorrentContent.InfoHash.EqCol(q.TorrentFile.InfoHash),
						},
						Type: query.TableJoinTypeInner,
					},
				}
			}),
		)
		if qErr != nil {
			errs = append(errs, qErr)
			return
		}
		if err := q.UnderlyingDB().Select(
			strings.Replace(ftFieldTemplate, "{table}", "torrent_files", -1),
			"count(distinct(torrent_files.info_hash)) as count",
		).Where("torrent_files.extension in " + makeStringList(allExts...)).Group(
			"file_type",
		).Find(&fileResults).Error; err != nil {
			errs = append(errs, err)
		}
	}()
	go func() {
		defer wg.Done()
		q, qErr := ctx.NewAggregationQuery(
			query.Table(model.TableNameTorrent),
			query.Join(func(q *dao.Query) []query.TableJoin {
				return []query.TableJoin{
					{
						Table: q.TorrentContent,
						On: []field.Expr{
							q.TorrentContent.InfoHash.EqCol(q.Torrent.InfoHash),
						},
						Type: query.TableJoinTypeInner,
					},
				}
			}),
		)
		if qErr != nil {
			errs = append(errs, qErr)
			return
		}
		if err := q.UnderlyingDB().Select(
			strings.Replace(ftFieldTemplate, "{table}", "torrents", -1),
			"count(*) as count",
		).Where("torrents.extension in " + makeStringList(allExts...)).Group(
			"file_type",
		).Find(&torrentResults).Error; err != nil {
			errs = append(errs, err)
		}
	}()
	wg.Wait()
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	allResults := make([]result, 0, len(fileResults)+len(torrentResults))
	allResults = append(allResults, fileResults...)
	allResults = append(allResults, torrentResults...)
	agg := make(query.AggregationItems, len(allResults))
	for _, item := range allResults {
		key := item.FileType.String()
		if existing, ok := agg[key]; !ok {
			agg[key] = query.AggregationItem{
				Label: item.FileType.Label(),
				Count: item.Count,
			}
		} else {
			existing.Count += item.Count
			agg[key] = existing
		}
	}
	return agg, nil
}

func makeStringList(values ...string) string {
	strs := "("
	for i, ext := range values {
		if i > 0 {
			strs += ","
		}
		strs += "'" + ext + "'"
	}
	strs += ")"
	return strs
}

func (f torrentFileTypeFacet) Criteria() []query.Criteria {
	return []query.Criteria{query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
		filter := f.Filter().Values()
		if len(filter) == 0 {
			return query.RawCriteria{
				Query: "1=1",
			}, nil
		}
		fileTypes := make([]model.FileType, 0, len(filter))
		for _, v := range filter {
			ft, ftErr := model.ParseFileType(v)
			if ftErr != nil {
				return nil, errors.New("invalid file type filter specified")
			}
			fileTypes = append(fileTypes, ft)
		}
		return TorrentFileTypeCriteria(fileTypes...), nil
	})}
}
