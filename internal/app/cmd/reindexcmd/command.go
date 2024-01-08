package reindexcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

type Params struct {
	fx.In
	Dao    lazy.Lazy[*dao.Query]
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	return Result{Command: &cli.Command{
		Name:  "reindex",
		Usage: "Reindex all records for full text search",
    Flags: []cli.Flag{
      &cli.IntFlag{
        Name: "batchSize",
        Value: 100,
      },
    },
		Action: func(ctx *cli.Context) error {
			println("reindexing...")
			d, err := p.Dao.Get()
			if err != nil {
				return err
			}
      batchSize := ctx.Int("batchSize")
			contentCount := int64(0)
			torrentContentCount := int64(0)
			if result, err := d.Content.WithContext(ctx.Context).Count(); err != nil {
				return err
			} else {
				contentCount = result
			}
			if result, err := d.TorrentContent.WithContext(ctx.Context).Count(); err != nil {
				return err
			} else {
				torrentContentCount = result
			}
			contentBar := progressbar.Default(contentCount, "[1/2] reindexing content")
			tsvs := make(map[model.ContentRef]fts.Tsvector)
			var contentResult []*model.Content
			if err := d.Content.WithContext(ctx.Context).Preload(
				d.Content.Attributes.RelationField,
				d.Content.Collections.RelationField,
			).FindInBatches(&contentResult, batchSize, func(tx gen.Dao, _ int) error {
				for _, c := range contentResult {
					c.UpdateTsv()
					c.Collections = nil
					c.Attributes = nil
					tsvs[c.Ref()] = c.Tsv
				}
				if err := tx.Clauses(
					clause.OnConflict{
						Columns:   []clause.Column{{Name: "type"}, {Name: "source"}, {Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"tsv", "updated_at"}),
					},
				).CreateInBatches(&contentResult, batchSize); err != nil {
					return err
				}
				_ = contentBar.Add(len(contentResult))
				return nil
			}); err != nil {
				return err
			}
			_ = contentBar.Finish()
			torrentContentBar := progressbar.Default(torrentContentCount, "[2/2] reindexing torrent content")
			var torrentContentResult []*model.TorrentContent
			if err := d.TorrentContent.WithContext(ctx.Context).Preload(
				d.TorrentContent.Torrent.RelationField,
				d.TorrentContent.Torrent.Files.RelationField,
			).FindInBatches(&torrentContentResult, batchSize, func(tx gen.Dao, _ int) error {
				for _, tc := range torrentContentResult {
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
				}
				if err := tx.Clauses(
					clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"tsv", "updated_at"}),
					},
				).CreateInBatches(torrentContentResult, batchSize); err != nil {
					return err
				}
				_ = torrentContentBar.Add(len(torrentContentResult))
				return nil
			}); err != nil {
				return err
			}
			_ = torrentContentBar.Finish()
			return nil
		},
	}}, nil
}
