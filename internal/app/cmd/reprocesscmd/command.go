package reprocesscmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen"
)

type Params struct {
	fx.In
	Dao                lazy.Lazy[*dao.Query]
	ProcessorPublisher lazy.Lazy[publisher.Publisher[processor.MessageParams]]
	Logger             *zap.SugaredLogger
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
				Value: 100,
			},
			&cli.StringFlag{
				Name:  "classifyMode",
				Value: "default",
				Usage: "default (only attempt to match previously unmatched torrents);\n" +
					"rematch (ignore any pre-existing classification and always classify from scratch);\n" +
					"skip (skip classification for previously unmatched torrents that don't have any hint)",
			},
		},
		Action: func(ctx *cli.Context) error {
			var classifyMode processor.ClassifyMode
			switch ctx.String("classifyMode") {
			case "default":
				classifyMode = processor.ClassifyModeDefault
			case "rematch":
				classifyMode = processor.ClassifyModeRematch
			case "skip":
				classifyMode = processor.ClassifyModeSkipUnmatched
			default:
				return cli.Exit("invalid classifyMode", 1)
			}
			println("queueing full reprocess...")
			d, err := p.Dao.Get()
			if err != nil {
				return err
			}
			p, err := p.ProcessorPublisher.Get()
			if err != nil {
				return err
			}
			batchSize := ctx.Int("batchSize")
			torrentCount := int64(0)
			if result, err := d.Torrent.WithContext(ctx.Context).Count(); err != nil {
				return err
			} else {
				torrentCount = result
			}
			bar := progressbar.Default(torrentCount, "queuing torrents")
			var torrentResult []*model.Torrent
			if err := d.Torrent.WithContext(ctx.Context).FindInBatches(&torrentResult, batchSize, func(tx gen.Dao, _ int) error {
				infoHashes := make([]protocol.ID, 0, len(torrentResult))
				for _, c := range torrentResult {
					infoHashes = append(infoHashes, c.InfoHash)
				}
				if _, err := p.Publish(ctx.Context, processor.MessageParams{
					ClassifyMode: classifyMode,
					InfoHashes:   infoHashes,
				}); err != nil {
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
