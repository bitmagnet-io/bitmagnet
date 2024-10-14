package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm/clause"
)

// TorrentFilesOrderBy represents sort orders for torrent content search results
// ENUM(index, path, extension, size)
type TorrentFilesOrderBy string

func (ob TorrentFilesOrderBy) Clauses(direction OrderDirection) []query.OrderByColumn {
	desc := direction == OrderDirectionDescending
	switch ob {
	case TorrentFilesOrderByIndex:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "index",
				},
				Desc: desc,
			},
		}}
	case TorrentFilesOrderByPath:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "path",
				},
				Desc: desc,
			},
		}}
	case TorrentFilesOrderByExtension:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameTorrentFile,
					Name:  "extension",
				},
				Desc: desc,
			},
		}}
	case TorrentFilesOrderBySize:
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

type TorrentFilesFullOrderBy maps.InsertMap[TorrentFilesOrderBy, OrderDirection]

func (fob TorrentFilesFullOrderBy) Clauses() []query.OrderByColumn {
	im := maps.InsertMap[TorrentFilesOrderBy, OrderDirection](fob)
	clauses := make([]query.OrderByColumn, 0, im.Len())
	for _, ob := range im.Entries() {
		clauses = append(clauses, ob.Key.Clauses(ob.Value)...)
	}
	return clauses
}

func (fob TorrentFilesFullOrderBy) Option() query.Option {
	return query.OrderBy(fob.Clauses()...)
}
