package httpserver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/settings"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Client   lazy.Lazy[torznab.Client]
	Settings lazy.Lazy[*settings.Settings]
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		Option: builder{
			client:   p.Client,
			settings: p.Settings,
		},
	}
}

type builder struct {
	client   lazy.Lazy[torznab.Client]
	settings lazy.Lazy[*settings.Settings]
}

func (builder) Key() string {
	return "torznab"
}

type torznabworker struct {
	client torznab.Client
	// settings *settings.Settings
	profile  settings.Profile
	hostname *string
}

func (w torznabworker) writeInternalError(c *gin.Context, err error) {
	_ = c.AbortWithError(500, err)
	_, _ = c.Writer.WriteString(err.Error() + "\n")
}

func (w torznabworker) writeXml(c *gin.Context, obj torznab.Xmler) {
	body, err := obj.Xml()
	if err != nil {
		w.writeInternalError(c, fmt.Errorf("failed to encode xml: %w", err))
		return
	}
	c.Status(200)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	_, _ = c.Writer.Write(body)
}

func (w torznabworker) writeErr(c *gin.Context, err error) {
	torznabErr := &torznab.Error{}
	if ok := errors.As(err, torznabErr); ok {
		w.writeXml(c, torznabErr)
	} else {
		w.writeInternalError(c, err)
	}
}

func (w torznabworker) permaLinkBase(c *gin.Context) string {
	if w.hostname != nil {
		return *w.hostname
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host + "/webui/torrents/permalink/"
}

func (w torznabworker) get(c *gin.Context) {
	tp := c.Query(torznab.ParamType)
	if tp == "" {
		w.writeErr(c, torznab.Error{
			Code:        200,
			Description: fmt.Sprintf("missing parameter (%s)", torznab.ParamType),
		})
		return
	}
	if tp == torznab.FunctionCaps {
		caps, capsErr := w.client.Caps(c)
		if capsErr != nil {
			w.writeErr(c, fmt.Errorf("failed to execute caps: %w", capsErr))
			return
		}
		w.writeXml(c, caps)
		return
	}
	var cats []int
	for _, csvCat := range c.QueryArray(torznab.ParamCat) {
		for _, cat := range strings.Split(csvCat, ",") {
			if intCat, err := strconv.Atoi(cat); err == nil {
				cats = append(cats, intCat)
			}
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
	result, searchErr := w.client.Search(c, torznab.SearchRequest{
		Query:          c.Query(torznab.ParamQuery),
		Type:           tp,
		Cats:           cats,
		ImdbId:         imdbId,
		TmdbId:         tmdbId,
		Season:         season,
		Episode:        episode,
		Limit:          limit,
		Offset:         offset,
		OrderBy:        w.profile.OrderBy,
		OrderDirection: w.profile.OrderDirection,
		Tags:           w.profile.Tags,
		PermaLinkBase:  w.permaLinkBase(c),
	})
	if searchErr != nil {
		w.writeErr(c, fmt.Errorf("failed to search: %w", searchErr))
		return
	}
	w.writeXml(c, result)

}

func (w torznabworker) getDefault(profile settings.Profile) gin.HandlerFunc {
	w.profile = profile
	return gin.HandlerFunc(w.get)
}

func (w torznabworker) getWithProfile(profiles map[string]settings.Profile) gin.HandlerFunc {
	handler := func(c *gin.Context) {
		profileName := c.Param("profile")
		profile, ok := profiles[profileName]
		if !ok {
			w.writeErr(c, torznab.Error{
				Code:        200,
				Description: fmt.Sprintf("profile not found (%s)", profileName),
			})
			return
		}
		w.profile = profile
		w.get(c)
	}
	return gin.HandlerFunc(handler)

}

func (b builder) Apply(e *gin.Engine) error {
	client, err := b.client.Get()
	if err != nil {
		return err
	}
	settings, err := b.settings.Get()
	if err != nil {
		return err
	}
	worker := torznabworker{
		client:   client,
		hostname: settings.Hostname,
	}

	e.GET("/torznab/api/*any", worker.getDefault(settings.Profiles[torznab.ProfileDefault]))
	e.GET("/torznab/:profile/*any", worker.getWithProfile(settings.Profiles))
	return nil
}
