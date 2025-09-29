package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"gorm.io/gorm/clause"
)

func TorrentContentOrderByClauses(ob adapter.TorrentContentOrderBy, desc bool) []query.OrderByColumn {
	switch ob {
	case adapter.TorrentContentOrderByRelevance:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: query.QueryStringRankField,
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderByPublishedAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "published_at",
				},
				Desc: desc,
			},
		}, {
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderByUpdatedAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "updated_at",
				},
				Desc: desc,
			},
		}, {
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderBySize:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "size",
				},
				Desc: desc,
			},
		}, {
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderByFilesCount:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: "COALESCE(" + model.TableNameTorrentContent + ".files_count, 0)",
					Raw:  true,
				},
				Desc: desc,
			},
		}, {
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderBySeeders:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: "coalesce(" + model.TableNameTorrentContent + ".seeders, -1)",
					Raw:  true,
				},
				Desc: desc,
			},
		}, {
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderByLeechers:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Name: "coalesce(" + model.TableNameTorrentContent + ".leechers, -1)",
					Raw:  true,
				},
				Desc: desc,
			},
		}, {
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentContentOrderByName:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrent,
					Name:  "name",
				},
				Desc: desc,
			},
			RequiredJoins: []string{model.TableNameTorrent},
		}}
	case adapter.TorrentContentOrderByInfoHash:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentContent,
					Name:  "info_hash",
				},
				Desc: desc,
			},
		}}
	default:
		return []query.OrderByColumn{}
	}
}
