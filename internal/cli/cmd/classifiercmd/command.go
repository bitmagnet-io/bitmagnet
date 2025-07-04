package classifiercmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/bitmagnet-io/bitmagnet/internal/app"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

var Command = &cli.Command{
	Name: "classifier",
	Commands: []*cli.Command{
		{
			Name:  "show",
			Usage: "Shows the classifier workflow source",
			Flags: []cli.Flag{
				&formatFlag,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fxApp := app.New(
					ctx,
					cmd.Writer,
					fx.Invoke(func(
						lWorkflowSource lazy.Lazy[classifier.Source],
					) error {
						workflowSource, err := lWorkflowSource.Get()
						if err != nil {
							return err
						}
						return write(cmd.Writer, workflowSource, cmd.String("format"))
					},
					),
				)

				fxApp.Run()

				return fxApp.Err()
			},
		},
		{
			Name:  "schema",
			Usage: "Shows the classifier JSON schema",
			Flags: []cli.Flag{
				&formatFlag,
			},
			Action: func(_ context.Context, cmd *cli.Command) error {
				return write(
					cmd.Writer,
					classifier.DefaultJSONSchema(),
					cmd.String("format"),
				)
			},
		},
	},
}

var formatFlag = cli.StringFlag{
	Name:  "format",
	Usage: "Output format (json or yaml)",
	Value: "yaml",
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
