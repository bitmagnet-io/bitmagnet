package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"gorm.io/gorm/clause"
)

func TorrentFilesOrderByClauses(ob adapter.TorrentFilesOrderBy, desc bool) []query.OrderByColumn {
	switch ob {
	case adapter.TorrentFilesOrderByIndex:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "index",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentFilesOrderByPath:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "path",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentFilesOrderByExtension:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "extension",
				},
				Desc: desc,
			},
		}}
	case adapter.TorrentFilesOrderBySize:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "size",
				},
				Desc: desc,
			},
		}}
	default:
		return []query.OrderByColumn{}
	}
}
