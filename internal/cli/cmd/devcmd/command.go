package devcmd

import (
	"github.com/bitmagnet-io/bitmagnet/internal/cli/cmd/devcmd/gormcmd"
	"github.com/bitmagnet-io/bitmagnet/internal/cli/cmd/devcmd/migratecmd"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "dev",
	Usage: "Development commands",
	Commands: []*cli.Command{
		gormcmd.Command,
		migratecmd.Command,
	},
}
