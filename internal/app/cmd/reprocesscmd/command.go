package reprocesscmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen"
	"strings"
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
		Name:  "reprocess",
		Usage: "Queue all torrents for reprocessing",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "batchSize",
				Value: 1000,
			},
			&cli.StringSliceFlag{
				Name:    "contentType",
				Aliases: []string{"contentTypes"},
				Usage: "reprocess only torrents with the specified content type(s) (e.g. '" +
					strings.Join(model.ContentTypeNames(), "', '") +
					"', or 'null' for unknown)",
			},
			&cli.BoolFlag{
				Name:  "orphans",
				Usage: "reprocess only torrents that have no torrent_contents record",
			},
			&cli.StringFlag{
				Name:  "classifyMode",
				Value: "default",
				Usage: "default (only attempt to match previously unmatched torrents);\n" +
					"rematch (ignore any pre-existing match and always classify from scratch)",
			},
		},
		Action: func(ctx *cli.Context) error {
			var classifyMode processor.ClassifyMode
			switch ctx.String("classifyMode") {
			case "default":
				classifyMode = processor.ClassifyModeDefault
			case "rematch":
				classifyMode = processor.ClassifyModeRematch
			default:
				return cli.Exit("invalid classifyMode", 1)
			}
			var contentTypes []string
			var unknownContentType bool
			for _, contentType := range ctx.StringSlice("contentType") {
				if contentType == "null" {
					unknownContentType = true
				} else {
					ct, err := model.ParseContentType(contentType)
					if err != nil {
						return err
					}
					contentTypes = append(contentTypes, ct.String())
				}
			}
			d, err := p.Dao.Get()
			if err != nil {
				return err
			}
			println("queueing reprocess...")
			var scopes []func(gen.Dao) gen.Dao
			if len(contentTypes) > 0 || unknownContentType {
				scopes = append(scopes, func(tx gen.Dao) gen.Dao {
					sq := d.TorrentContent.Where(
						d.TorrentContent.InfoHash.EqCol(d.Torrent.InfoHash),
					).Where(d.TorrentContent.ContentType.In(contentTypes...))
					if unknownContentType {
						sq = sq.Or(d.TorrentContent.ContentType.IsNull())
					}
					return tx.Where(gen.Exists(sq))
				})
			}
			if ctx.Bool("orphans") {
				scopes = append(scopes, func(tx gen.Dao) gen.Dao {
					return tx.Not(
						gen.Exists(
							d.TorrentContent.Where(
								d.TorrentContent.InfoHash.EqCol(d.Torrent.InfoHash),
							),
						),
					)
				})
			}
			batchSize := ctx.Int("batchSize")
			torrentCount := int64(0)
			if result, err := dao.BudgetedCount(d.Torrent.WithContext(ctx.Context).Scopes(scopes...).UnderlyingDB(), 10_000); err != nil {
				return err
			} else {
				torrentCount = result.Count
			}
			bar := progressbar.Default(torrentCount, "queuing torrents")
			var torrentResult []*model.Torrent
			if err := d.Torrent.WithContext(ctx.Context).Scopes(scopes...).FindInBatches(&torrentResult, batchSize, func(tx gen.Dao, _ int) error {
				infoHashes := make([]protocol.ID, 0, len(torrentResult))
				for _, c := range torrentResult {
					infoHashes = append(infoHashes, c.InfoHash)
				}
				job, err := processor.NewQueueJob(processor.MessageParams{
					ClassifyMode: classifyMode,
					InfoHashes:   infoHashes,
				}, model.QueueJobPriority(10))
				if err != nil {
					return err
				}
				if err := tx.Create(&job); err != nil {
					return err
				}
				_ = bar.Add(len(torrentResult))
				return nil
			}); err != nil {
				return err
			}
			_ = bar.Finish()
			return nil
		},
	}}, nil
}
