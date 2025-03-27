package args

import (
	"os"

	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Args []string `name:"cli_args"`
}

func New() (Result, error) {
	return Result{Args: os.Args}, nil
}
