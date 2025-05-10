package classifiercmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

type Params struct {
	fx.In
	WorkflowSource lazy.Lazy[classifier.Source]
}

type Result struct {
	fx.Out
	Command *cli.Command `group:"commands"`
}

var formatFlag = cli.StringFlag{
	Name:  "format",
	Usage: "Output format (json or yaml)",
	Value: "yaml",
}

func New(p Params) (Result, error) {
	return Result{Command: &cli.Command{
		Name: "classifier",
		Subcommands: []*cli.Command{
			{
				Name:  "show",
				Usage: "Show the classifier workflow source",
				Flags: []cli.Flag{
					&formatFlag,
				},
				Action: func(ctx *cli.Context) error {
					src, srcErr := p.WorkflowSource.Get()
					if srcErr != nil {
						return srcErr
					}
					return write(ctx.App.Writer, src, ctx.String("format"))
				},
			},
			{
				Name:  "schema",
				Usage: "Show the classifier JSON schema",
				Flags: []cli.Flag{
					&formatFlag,
				},
				Action: func(ctx *cli.Context) error {
					return write(
						ctx.App.Writer,
						classifier.DefaultJSONSchema(),
						ctx.String("format"),
					)
				},
			},
		},
	}}, nil
}

func write(writer io.Writer, src any, format string) error {
	var (
		output    []byte
		outputErr error
	)

	switch format {
	case "json":
		output, outputErr = json.MarshalIndent(src, "", "  ")
		output = append(output, '\n')
	case "yaml":
		output, outputErr = yaml.Marshal(src)
	default:
		outputErr = fmt.Errorf("unsupported format: %s", format)
	}

	if outputErr != nil {
		return outputErr
	}

	_, writeErr := writer.Write(output)

	return writeErr
}
