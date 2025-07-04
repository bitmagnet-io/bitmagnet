package cli

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/cli/cmd/classifiercmd"
	"github.com/bitmagnet-io/bitmagnet/internal/cli/cmd/configcmd"
	"github.com/bitmagnet-io/bitmagnet/internal/cli/cmd/workercmd"

	"github.com/bitmagnet-io/bitmagnet/internal/colors"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/urfave/cli/v3"
)

var App = &cli.Command{
	Name:  "bitmagnet",
	Usage: "BitTorrent indexer and DHT crawler",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-banner",
			Usage: "disable banner",
		},
	},
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		if !cmd.Bool("no-banner") {
			writeBanner(cmd.Writer)
		}

		return ctx, nil
	},
	Commands: []*cli.Command{
		classifiercmd.Command,
		configcmd.Command,
		workercmd.Command,
	},
	Version:         version.GitTag,
	HideHelpCommand: true,
	ExitErrHandler: func(_ context.Context, cmd *cli.Command, err error) {
		_, _ = fmt.Fprintf(
			cmd.ErrWriter,
			"\n%s %s\n",
			colors.Red.Sprintf("execution failed:"),
			err,
		)
	},
}
