package torrentcmd

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm/clause"
	"net/netip"
)

type Params struct {
	fx.In
	MetaInfoRequester metainforequester.Requester
	Classifier        classifier.Classifier
	Dao               *dao.Query
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
					infoHash, err := protocol.ParseID(ctx.String("infoHash"))
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
			{
				Name: "updateAll",
				Action: func(ctx *cli.Context) error {
					var tcs []*model.TorrentContent
					result := p.Dao.WithContext(ctx.Context).TorrentContent.Preload(
						p.Dao.TorrentContent.Torrent.RelationField,
						//p.Dao.TorrentContent.Content.RelationField,
						//p.Dao.TorrentContent.Content.Attributes.RelationField,
						//p.Dao.TorrentContent.Content.Collections.RelationField,
					).Where(p.Dao.TorrentContent.ContentID.IsNull()).FindInBatches(&tcs, 1000, func(tx gen.Dao, batch int) error {
						for _, tc := range tcs {
							//fmt.Printf("a: %#v\n", tc)
							if err := tc.UpdateFields(); err != nil {
								return err
							}
							tc.Torrent = model.Torrent{}
							tc.Content = model.Content{}
							//fmt.Printf("b: %#v\n", tc)
						}

						// Save changes to the records in the current batch
						return tx.Clauses(clause.OnConflict{
							UpdateAll: true,
						}).Save(&tcs)
					})
					fmt.Printf("result: %v\n", result)
					return nil
				},
			},
		},
	}}, nil
}
