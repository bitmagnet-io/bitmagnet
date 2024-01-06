package torrentcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
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
	Classifier        lazy.Lazy[classifier.Classifier]
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
				Name: "classify",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "infoHash",
					},
				},
				Action: func(ctx *cli.Context) error {
					c, err := p.Classifier.Get()
					if err != nil {
						return err
					}
					infoHash, err := protocol.ParseID(ctx.String("infoHash"))
					if err != nil {
						return err
					}
					return c.Classify(ctx.Context, infoHash)
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
