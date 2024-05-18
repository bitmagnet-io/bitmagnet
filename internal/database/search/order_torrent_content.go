package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm/clause"
)

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -f order_torrent_content.go

// TorrentContentOrderBy represents sort orders for torrent content search results
// ENUM(Relevance, PublishedAt, UpdatedAt, Size, Files, Seeders, Leechers, Name, InfoHash)
type TorrentContentOrderBy string

// OrderDirection represents sort order directions
// ENUM(Ascending, Descending)
type OrderDirection string

func (ob TorrentContentOrderBy) Clauses(direction OrderDirection) []clause.OrderByColumn {
	desc := direction == OrderDirectionDescending
	switch ob {
	case TorrentContentOrderByRelevance:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Name: query.QueryStringRankField,
			},
			Desc: desc,
		}}
	case TorrentContentOrderByPublishedAt:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Table: model.TableNameTorrentContent,
				Name:  "published_at",
			},
			Desc: desc,
		}}
	case TorrentContentOrderByUpdatedAt:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Table: model.TableNameTorrentContent,
				Name:  "updated_at",
			},
			Desc: desc,
		}}
	case TorrentContentOrderBySize:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Table: model.TableNameTorrent,
				Name:  "size",
			},
			Desc: desc,
		}}
	case TorrentContentOrderByFiles:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Name: "CASE WHEN " + model.TableNameTorrent + ".files_status = 'single' THEN 1 ELSE COALESCE(" + model.TableNameTorrent + ".files_count, -1) END",
				Raw:  true,
			},
			Desc: desc,
		}}
	case TorrentContentOrderBySeeders:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Name: "COALESCE(" + model.TableNameTorrentContent + ".seeders, -1)",
				Raw:  true,
			},
			Desc: desc,
		}}
	case TorrentContentOrderByLeechers:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Name: "COALESCE(" + model.TableNameTorrentContent + ".leechers, -1)",
				Raw:  true,
			},
			Desc: desc,
		}}
	case TorrentContentOrderByName:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Table: model.TableNameTorrent,
				Name:  "name",
			},
			Desc: desc,
		}}
	case TorrentContentOrderByInfoHash:
		return []clause.OrderByColumn{{
			Column: clause.Column{
				Table: model.TableNameTorrentContent,
				Name:  "info_hash",
			},
			Desc: desc,
		}}
	default:
		return []clause.OrderByColumn{}
	}
}

type TorrentContentFullOrderBy maps.InsertMap[TorrentContentOrderBy, OrderDirection]

func (fob TorrentContentFullOrderBy) Clauses() []clause.OrderByColumn {
	im := maps.InsertMap[TorrentContentOrderBy, OrderDirection](fob)
	clauses := make([]clause.OrderByColumn, 0, im.Len()+2)
	for _, ob := range im.Entries() {
		clauses = append(clauses, ob.Key.Clauses(ob.Value)...)
	}
	// make ordering alphabetical and deterministic when not already specified:
	if !im.Has(TorrentContentOrderByName) {
		clauses = append(clauses, TorrentContentOrderByName.Clauses(OrderDirectionAscending)...)
	}
	if !im.Has(TorrentContentOrderByInfoHash) {
		clauses = append(clauses, TorrentContentOrderByInfoHash.Clauses(OrderDirectionAscending)...)
	}
	return clauses
}

func (fob TorrentContentFullOrderBy) Option() query.Option {
	return query.Options(
		query.RequireJoin(model.TableNameTorrent),
		query.OrderBy(fob.Clauses()...),
	)
}
