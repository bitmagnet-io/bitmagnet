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

func (ob TorrentContentOrderBy) Clauses(direction OrderDirection) []query.OrderByColumn {
	desc := direction == OrderDirectionDescending
	switch ob {
	case TorrentContentOrderByRelevance:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: query.QueryStringRankField,
				},
				Desc: desc,
			},
		}}
	case TorrentContentOrderByPublishedAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "published_at",
				},
				Desc: desc,
			}}}
	case TorrentContentOrderByUpdatedAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "updated_at",
				},
				Desc: desc,
			}}}
	case TorrentContentOrderBySize:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "size",
				},
				Desc: desc,
			}}}
	case TorrentContentOrderByFiles:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: "COALESCE(" + model.TableNameTorrentContent + ".files_count, 0)",
					Raw:  true,
				},
				Desc: desc,
			}}}
	case TorrentContentOrderBySeeders:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: "coalesce(" + model.TableNameTorrentContent + ".seeders, -1)",
					Raw:  true,
				},
				Desc: desc,
			},
		}}
	case TorrentContentOrderByLeechers:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: "coalesce(" + model.TableNameTorrentContent + ".leechers, -1)",
					Raw:  true,
				},
				Desc: desc,
			},
		}}
	case TorrentContentOrderByName:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrent,
					Name:  "name",
				},
				Desc: desc,
			},
			RequiredJoins: []string{model.TableNameTorrent}}}
	case TorrentContentOrderByInfoHash:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			}}}
	default:
		return []query.OrderByColumn{}
	}
}

type TorrentContentFullOrderBy maps.InsertMap[TorrentContentOrderBy, OrderDirection]

func (fob TorrentContentFullOrderBy) Clauses() []query.OrderByColumn {
	im := maps.InsertMap[TorrentContentOrderBy, OrderDirection](fob)
	clauses := make([]query.OrderByColumn, 0, im.Len())
	for _, ob := range im.Entries() {
		clauses = append(clauses, ob.Key.Clauses(ob.Value)...)
	}
	return clauses
}

func (fob TorrentContentFullOrderBy) Option() query.Option {
	return query.OrderBy(fob.Clauses()...)
}
