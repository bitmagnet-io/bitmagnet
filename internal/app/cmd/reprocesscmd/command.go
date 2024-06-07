package reprocesscmd

import (
	"encoding/json"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
			&cli.UintFlag{
				Name:  "batchSize",
				Value: 100,
			},
			&cli.UintFlag{
				Name:  "chunkSize",
				Value: 10_000,
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
			&cli.StringFlag{
				Name:  "classifierFlags",
				Value: "{}",
				Usage: "optional JSON-encoded runtime flags to pass to the classifier",
			},
			&cli.BoolFlag{
				Name:  "apisDisabled",
				Value: false,
				Usage: "disable API calls for the classifier workflow",
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
			var flags classifier.Flags
			strFlags := ctx.String("classifierFlags")
			if err := json.Unmarshal([]byte(strFlags), &flags); err != nil {
				return cli.Exit("invalid flags", 1)
			}
			if ctx.Bool("apisDisabled") {
				flags["apis_enabled"] = false
			}
			var contentTypes []model.NullContentType
			for _, contentType := range ctx.StringSlice("contentType") {
				if contentType == "null" {
					contentTypes = append(contentTypes, model.NullContentType{})
				} else {
					ct, err := model.ParseContentType(contentType)
					if err != nil {
						return err
					}
					contentTypes = append(contentTypes, model.NullContentType{
						ContentType: ct,
						Valid:       true,
					})
				}
			}
			job, err := batch.NewQueueJob(batch.MessageParams{
				ClassifyMode:    classifyMode,
				ClassifierFlags: flags,
				ChunkSize:       ctx.Uint("chunkSize"),
				BatchSize:       ctx.Uint("batchSize"),
				ContentTypes:    contentTypes,
				Orphans:         ctx.Bool("orphans"),
			})
			if err != nil {
				return err
			}
			d, err := p.Dao.Get()
			if err != nil {
				return err
			}
			createErr := d.QueueJob.WithContext(ctx.Context).Create(&job)
			var pgErr *pgconn.PgError
			if errors.As(createErr, &pgErr) && pgErr.Code == "23505" {
				_, _ = ctx.App.ErrWriter.Write([]byte("Reprocess already queued!\n"))
			} else if err == nil {
				_, _ = ctx.App.Writer.Write([]byte("Reprocess queued!\n"))
			}
			return createErr
		},
	}}, nil
}
