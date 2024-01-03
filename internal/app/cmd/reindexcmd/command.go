package reindexcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type Params struct {
	fx.In
	Search search.Search
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	return Result{Command: &cli.Command{
		Name: "reindex",
		Action: func(ctx *cli.Context) error {
			println("reindexing...")
			contentCount := int64(0)
			torrentContentCount := int64(0)
			if result, err := p.Search.Content(ctx.Context, query.WithTotalCount(true), query.Limit(0)); err != nil {
				return err
			} else {
				contentCount = int64(result.TotalCount)
			}
			if result, err := p.Search.TorrentContent(ctx.Context, query.WithTotalCount(true), query.Limit(0)); err != nil {
				return err
			} else {
				torrentContentCount = int64(result.TotalCount)
			}
			contentBar := progressbar.Default(contentCount, "[1/2] reindexing content")
			tsvs := make(map[model.ContentRef]fts.Tsvector)
			if err := p.Search.ContentBatch(ctx.Context, func(tx *dao.Query, r []search.ContentResultItem) error {
				cs := make([]*model.Content, 0, len(r))
				for _, item := range r {
					c := item.Content
					c.UpdateTsv()
					c.Collections = nil
					c.Attributes = nil
					cs = append(cs, &c)
					tsvs[c.Ref()] = c.Tsv
				}
				if err := tx.Content.WithContext(ctx.Context).Clauses(
					clause.OnConflict{
						Columns:   []clause.Column{{Name: "type"}, {Name: "source"}, {Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"tsv", "updated_at"}),
					},
				).CreateInBatches(cs, 1000); err != nil {
					return err
				}
				_ = contentBar.Add(len(cs))
				return nil
			},
				search.ContentCoreJoins(),
				search.ContentDefaultPreload(),
				search.ContentDefaultHydrate(),
				query.OrderBy(
					clause.OrderByColumn{Column: clause.Column{Name: "content.type"}},
					clause.OrderByColumn{Column: clause.Column{Name: "content.source"}},
					clause.OrderByColumn{Column: clause.Column{Name: "content.id"}},
				),
			); err != nil {
				return err
			}
			_ = contentBar.Finish()
			println(len(tsvs))
			torrentContentBar := progressbar.Default(torrentContentCount, "[2/2] reindexing torrent content")
			if err := p.Search.TorrentContentBatch(ctx.Context, func(tx *dao.Query, r []search.TorrentContentResultItem) error {
				tcs := make([]*model.TorrentContent, 0, len(r))
				for _, item := range r {
					tc := item.TorrentContent
					ref := tc.EntityReference()
					if ref.Valid {
						tsv, ok := tsvs[ref.Val]
						if !ok {
							p.Logger.Warnw("missing tsv", "ref", ref.Val)
							continue
						} else {
							tc.Content.Tsv = tsv
						}
					}
					tc.UpdateTsv()
					tc.Torrent = model.Torrent{}
					tc.Content = model.Content{}
					tcs = append(tcs, &tc)
				}
				if err := tx.TorrentContent.WithContext(ctx.Context).Clauses(
					clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"tsv", "updated_at"}),
					},
				).CreateInBatches(tcs, 1000); err != nil {
					return err
				}
				_ = torrentContentBar.Add(len(tcs))
				return nil
			},
				search.TorrentContentCoreJoins(),
				search.HydrateTorrentContentTorrent(),
				//search.TorrentContentDefaultHydrate(),
				query.OrderByColumn("torrent_contents.id", false),
			); err != nil {
				return err
			}
			_ = torrentContentBar.Finish()
			return nil
		},
	}}, nil
}
