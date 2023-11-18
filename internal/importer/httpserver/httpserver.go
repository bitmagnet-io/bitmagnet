package httpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/importer"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Params struct {
	fx.In
	Importer importer.Importer
	Logger   *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) (r Result, err error) {
	r.Option = &builder{p.handler}
	return
}

const ImportIdHeader = "x-import-id"

func (p Params) handler(ctx *gin.Context) {
	s := bufio.NewScanner(ctx.Request.Body)
	s.Split(bufio.ScanRunes)
	importId := ctx.Request.Header.Get(ImportIdHeader)
	if importId == "" {
		importId = strconv.FormatUint(uint64(time.Now().Unix()), 10)
	}
	i := p.Importer.New(ctx, importer.Info{
		ID: importId,
	})
	var currentLine []rune
	count := 0
	writeCount := func() {
		_, _ = ctx.Writer.WriteString(fmt.Sprintf("%d items imported\n", count))
	}
	addItem := func() error {
		item := importer.Item{}
		if err := json.Unmarshal([]byte(string(currentLine)), &item); err != nil {
			p.Logger.Errorw("error adding item", "error", err)
			ctx.Status(400)
			_, _ = ctx.Writer.WriteString(err.Error())
			return err
		}
		if err := i.Import(item); err != nil {
			p.Logger.Errorw("error importing item", "error", err)
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
			if ch == '\n' {
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
	i.Drain()
	if err := i.Close(); err != nil {
		p.Logger.Errorw("error closing import", "error", err)
		ctx.Status(400)
		_, _ = ctx.Writer.WriteString(err.Error())
		return
	}
	ctx.Status(200)
	writeCount()
}

type builder struct {
	handler gin.HandlerFunc
}

func (builder) Key() string {
	return "import"
}

func (b builder) Apply(e *gin.Engine) error {
	e.POST("/import", b.handler)
	return nil
}
