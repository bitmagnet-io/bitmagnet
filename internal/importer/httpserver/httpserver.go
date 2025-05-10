package httpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/importer"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Importer lazy.Lazy[importer.Importer]
	Logger   *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		Option: &builder{
			importer: p.Importer,
			logger:   p.Logger.Named("importer"),
		},
	}
}

const ImportIDHeader = "X-Import-Id"

type builder struct {
	importer lazy.Lazy[importer.Importer]
	logger   *zap.SugaredLogger
}

func (builder) Key() string {
	return "import"
}

func (b builder) Apply(e *gin.Engine) error {
	i, err := b.importer.Get()
	if err != nil {
		return err
	}

	e.POST("/import", func(ctx *gin.Context) {
		b.handle(ctx, i)
	})

	return nil
}

func (b builder) handle(ctx *gin.Context, i importer.Importer) {
	s := bufio.NewScanner(ctx.Request.Body)
	s.Split(bufio.ScanRunes)

	importID := ctx.Request.Header.Get(ImportIDHeader)
	if importID == "" {
		importID = strconv.FormatUint(uint64(time.Now().Unix()), 10)
	}

	ai := i.New(ctx, importer.Info{
		ID: importID,
	})

	var currentLine []rune

	count := 0
	writeCount := func() {
		_, _ = fmt.Fprintf(ctx.Writer, "%d items imported\n", count)
	}
	addItem := func() error {
		item := importer.Item{}
		// todo: Fix this
		//nolint:musttag
		if err := json.Unmarshal([]byte(string(currentLine)), &item); err != nil {
			b.logger.Errorw("error adding item", "error", err)
			ctx.Status(400)
			_, _ = ctx.Writer.WriteString(err.Error())

			return err
		}

		if err := ai.Import(item); err != nil {
			b.logger.Errorw("error importing item", "error", err)
			ctx.Status(400)
			_, _ = ctx.Writer.WriteString(err.Error())

			return err
		}

		count++
		if count%1_000 == 0 {
			writeCount()

			if count%10_000 == 0 {
				ctx.Writer.Flush()
			}
		}

		return nil
	}

	for s.Scan() {
		for _, ch := range s.Text() {
			if ch == '\n' && len(currentLine) > 0 {
				if err := addItem(); err != nil {
					return
				}

				currentLine = nil
			} else {
				currentLine = append(currentLine, ch)
			}
		}
	}

	if len(currentLine) > 0 {
		if err := addItem(); err != nil {
			return
		}
	}

	ai.Drain()

	if err := ai.Close(); err != nil {
		b.logger.Errorw("error closing import", "error", err)
		ctx.Status(400)
		_, _ = ctx.Writer.WriteString(err.Error())

		return
	}

	ctx.Status(200)
	writeCount()

	_, _ = ctx.Writer.WriteString("import complete\n")
}
