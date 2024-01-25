package torrentcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/netip"
)

type Params struct {
	fx.In
	MetaInfoRequester metainforequester.Requester
	Processor         lazy.Lazy[processor.Processor]
	Logger            *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

func New(p Params) (Result, error) {
	return Result{Command: &cli.Command{
		Name: "torrent",
		Subcommands: []*cli.Command{
			{
				Name: "process",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name: "infoHash",
					},
				},
				Action: func(ctx *cli.Context) error {
					pr, err := p.Processor.Get()
					if err != nil {
						return err
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
						ClassifyMode: processor.ClassifyModeRematch,
						InfoHashes:   infoHashes,
					})
				},
			},
			{
				Name: "requestMetaInfo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "infoHash",
					},
					&cli.StringFlag{
						Name: "address",
					},
				},
				Action: func(ctx *cli.Context) error {
					infoHash, err := protocol.ParseID(ctx.String("infoHash"))
					if err != nil {
						return err
					}
					addr, err := netip.ParseAddrPort(ctx.String("address"))
					if err != nil {
						return err
					}
					info, err := p.MetaInfoRequester.Request(ctx.Context, protocol.ID(infoHash), addr)
					if err != nil {
						return err
					}
					p.Logger.Infow("got infoBytes", "info", info)
					return nil
				},
			},
		},
	}}, nil
}
