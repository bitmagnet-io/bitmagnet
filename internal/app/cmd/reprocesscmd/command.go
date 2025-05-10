package reprocesscmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
			&cli.BoolFlag{
				Name:  "localSearchDisabled",
				Value: false,
				Usage: "disable local search queries for the classifier workflow",
			},
		},
		Action: p.action,
	}}, nil
}

func (p Params) action(ctx *cli.Context) error {
	var classifyMode processor.ClassifyMode

	switch ctx.String("classifyMode") {
	case "default":
		classifyMode = processor.ClassifyModeDefault
	case "rematch":
		classifyMode = processor.ClassifyModeRematch
	default:
		return errors.New("invalid classifyMode")
	}

	var flags classifier.Flags

	strFlags := ctx.String("classifierFlags")
	if err := json.Unmarshal([]byte(strFlags), &flags); err != nil {
		return fmt.Errorf("invalid flags: %w", err)
	}

	if ctx.Bool("apisDisabled") {
		flags["apis_enabled"] = false
	}

	if ctx.Bool("localSearchDisabled") {
		flags["local_search_enabled"] = false
	}

	contentTypes, err := getContentTypes(ctx)
	if err != nil {
		return err
	}

	job, err := batch.NewQueueJob(batch.MessageParams{
		ClassifyMode:    classifyMode,
		ClassifierFlags: flags,
		ChunkSize:       ctx.Uint("chunkSize"),
		BatchSize:       ctx.Uint("batchSize"),
		ContentTypes:    contentTypes,
		Orphans:         ctx.Bool("orphans"),
		UpdatedBefore:   time.Now(),
	})
	if err != nil {
		return err
	}

	d, err := p.Dao.Get()
	if err != nil {
		return err
	}

	if err := d.QueueJob.WithContext(ctx.Context).Create(&job); err != nil {
		return err
	}

	_, _ = ctx.App.Writer.Write([]byte("Reprocess queued!\n"))

	return nil
}

func getContentTypes(ctx *cli.Context) ([]model.NullContentType, error) {
	var contentTypes []model.NullContentType

	for _, contentType := range ctx.StringSlice("contentType") {
		if contentType == "null" {
			contentTypes = append(contentTypes, model.NullContentType{})
		} else {
			ct, err := model.ParseContentType(contentType)
			if err != nil {
				return nil, err
			}

			contentTypes = append(contentTypes, model.NullContentType{
				ContentType: ct,
				Valid:       true,
			})
		}
	}

	return contentTypes, nil
}
