//go:build dev

package cli

import (
	"github.com/bitmagnet-io/bitmagnet/internal/cli/cmd/devcmd"
)

func init() {
	App.Commands = append(App.Commands, devcmd.Command)
}
