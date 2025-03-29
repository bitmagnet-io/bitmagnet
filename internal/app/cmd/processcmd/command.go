package processcmd

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Processor lazy.Lazy[processor.Processor]
	Logger    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	return Result{
		Command: &cli.Command{
			Name: "process",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name: "infoHash",
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
			Action: func(ctx *cli.Context) error {
				pr, err := p.Processor.Get()
				if err != nil {
					return err
				}
				var flags classifier.Flags
				strFlags := ctx.String("classifierFlags")
				if err := json.Unmarshal([]byte(strFlags), &flags); err != nil {
					return cli.Exit("invalid flags", 1)
				}
				if ctx.Bool("apisDisabled") {
					flags["apis_enabled"] = false
				}
				if ctx.Bool("localSearchDisabled") {
					flags["local_search_enabled"] = false
				}
				var infoHashes []protocol.ID
				for _, infoHash := range ctx.StringSlice("infoHash") {
					id, err := protocol.ParseID(infoHash)
					if err != nil {
						return err
					}
					infoHashes = append(infoHashes, id)
				}
				if err != nil {
					return err
				}
				return pr.Process(ctx.Context, processor.MessageParams{
					ClassifyMode:    processor.ClassifyModeRematch,
					ClassifierFlags: flags,
					InfoHashes:      infoHashes,
				})
			},
		},
	}, nil
}
