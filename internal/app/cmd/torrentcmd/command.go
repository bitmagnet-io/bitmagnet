package torrentcmd

import (
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/netip"
)

type Params struct {
	fx.In
	MetaInfoRequester metainforequester.Requester
	Classifier        classifier.Classifier
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
					infoHash, err := model.NewHash20FromString(ctx.String("infoHash"))
					if err != nil {
						return err
					}
					return p.Classifier.Classify(ctx.Context, infoHash)
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
					infoHash, err := model.NewHash20FromString(ctx.String("infoHash"))
					if err != nil {
						return err
					}
					addr, err := netip.ParseAddrPort(ctx.String("address"))
					if err != nil {
						return err
					}
					nodeAddr := krpc.NodeAddr{}
					nodeAddr.FromAddrPort(addr)
					info, err := p.MetaInfoRequester.Request(ctx.Context, infoHash, nodeAddr)
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
