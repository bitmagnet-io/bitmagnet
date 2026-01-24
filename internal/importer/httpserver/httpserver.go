package httpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/importer"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(importer importer.Importer, logger *zap.Logger) gin.OptionFunc {
	return builder{
		importer: importer,
		logger:   logger.Named("importer"),
	}.Apply
}

const ImportIDHeader = "X-Import-Id"

type builder struct {
	importer importer.Importer
	logger   *zap.Logger
}

func (b builder) Apply(e *gin.Engine) {
	e.POST("/import", func(ctx *gin.Context) {
		b.handle(ctx, b.importer)
	})
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
		item := importer.Torrent{}

		//nolint:musttag
		if err := json.Unmarshal([]byte(string(currentLine)), &item); err != nil {
			b.logger.Error("failed to add item", zap.Error(err))
			ctx.Status(400)
			_, _ = ctx.Writer.WriteString(err.Error())

			return err
		}

		if err := ai.Import(item); err != nil {
			b.logger.Error("feiled to import item", zap.Error(err))
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
		b.logger.Error("failed to close import", zap.Error(err))
		ctx.Status(400)
		_, _ = ctx.Writer.WriteString(err.Error())

		return
	}

	ctx.Status(200)
	writeCount()

	_, _ = ctx.Writer.WriteString("import complete\n")
}
