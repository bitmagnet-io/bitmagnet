package httpserver

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"strconv"
)

type Params struct {
	fx.In
	Client lazy.Lazy[torznab.Client]
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		Option: builder{
			client: p.Client,
		},
	}
}

type builder struct {
	client lazy.Lazy[torznab.Client]
}

func (builder) Key() string {
	return "torznab"
}

func (b builder) Apply(e *gin.Engine) error {
	client, err := b.client.Get()
	if err != nil {
		return err
	}
	e.GET("/torznab/*any", func(c *gin.Context) {
		writeInternalError := func(err error) {
			_ = c.AbortWithError(500, err)
			_, _ = c.Writer.WriteString(err.Error() + "\n")
		}
		writeXml := func(obj torznab.Xmler) {
			body, err := obj.Xml()
			if err != nil {
				writeInternalError(fmt.Errorf("failed to encode xml: %w", err))
				return
			}
			c.Status(200)
			c.Header("Content-Type", "application/xml; charset=utf-8")
			_, _ = c.Writer.Write(body)
		}
		writeErr := func(err error) {
			torznabErr := &torznab.Error{}
			if ok := errors.As(err, torznabErr); ok {
				writeXml(torznabErr)
			} else {
				writeInternalError(err)
			}
		}
		tp := c.Query(torznab.ParamType)
		if tp == "" {
			writeErr(torznab.Error{
				Code:        200,
				Description: fmt.Sprintf("missing parameter (%s)", torznab.ParamType),
			})
			return
		}
		if tp == torznab.FunctionCaps {
			caps, capsErr := client.Caps(c)
			if capsErr != nil {
				writeErr(fmt.Errorf("failed to execute caps: %w", capsErr))
				return
			}
			writeXml(caps)
			return
		}
		var cats []int
		for _, cat := range c.QueryArray(torznab.ParamCat) {
			if intCat, err := strconv.Atoi(cat); err == nil {
				cats = append(cats, intCat)
			}
		}
		imdbId := model.NullString{}
		if qImdbId := c.Query(torznab.ParamImdbId); qImdbId != "" {
			imdbId.Valid = true
			imdbId.String = qImdbId
		}
		tmdbId := model.NullString{}
		if qTmdbId := c.Query(torznab.ParamTmdbId); qTmdbId != "" {
			tmdbId.Valid = true
			tmdbId.String = qTmdbId
		}
		season := model.NullInt{}
		episode := model.NullInt{}
		if qSeason := c.Query(torznab.ParamSeason); qSeason != "" {
			if intSeason, err := strconv.Atoi(qSeason); err == nil {
				season.Valid = true
				season.Int = intSeason
				if qEpisode := c.Query(torznab.ParamEpisode); qEpisode != "" {
					if intEpisode, err := strconv.Atoi(qEpisode); err == nil {
						episode.Valid = true
						episode.Int = intEpisode
					}
				}
			}
		}
		limit := model.NullUint{}
		if intLimit, limitErr := strconv.Atoi(c.Query(torznab.ParamLimit)); limitErr == nil && intLimit > 0 {
			limit.Valid = true
			limit.Uint = uint(intLimit)
		}
		offset := model.NullUint{}
		if intOffset, offsetErr := strconv.Atoi(c.Query(torznab.ParamOffset)); offsetErr == nil {
			offset.Valid = true
			offset.Uint = uint(intOffset)
		}
		result, searchErr := client.Search(c, torznab.SearchRequest{
			Query:   c.Query(torznab.ParamQuery),
			Type:    tp,
			Cats:    cats,
			ImdbId:  imdbId,
			TmdbId:  tmdbId,
			Season:  season,
			Episode: episode,
			Limit:   limit,
			Offset:  offset,
		})
		if searchErr != nil {
			writeErr(fmt.Errorf("failed to search: %w", searchErr))
			return
		}
		writeXml(result)
	})
	return nil
}
