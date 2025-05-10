package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/gin-gonic/gin"
)

type handler struct {
	config torznab.Config
	client torznab.Client
}

func (h handler) handleRequest(ctx *gin.Context) {
	profile, err := h.getProfile(ctx)
	if err != nil {
		h.writeError(ctx, err)
		return
	}

	tp := ctx.Query(torznab.ParamType)

	switch tp {
	case "":
		h.writeError(ctx, torznab.Error{
			Code:        200,
			Description: fmt.Sprintf("missing parameter (%s)", torznab.ParamType),
		})

	case torznab.FunctionCaps:
		h.writeXML(ctx, profile.Caps())

	default:
		h.handleSearch(ctx, profile, tp)
	}
}

func (h handler) handleSearch(ctx *gin.Context, profile torznab.Profile, tp string) {
	var cats []int

	for _, csvCat := range ctx.QueryArray(torznab.ParamCat) {
		for _, cat := range strings.Split(csvCat, ",") {
			if intCat, err := strconv.Atoi(cat); err == nil {
				cats = append(cats, intCat)
			}
		}
	}

	imdbID := model.NullString{}
	if qIMDBID := ctx.Query(torznab.ParamIMDBID); qIMDBID != "" {
		imdbID.Valid = true
		imdbID.String = qIMDBID
	}

	tmdbID := model.NullString{}
	if qTMDBID := ctx.Query(torznab.ParamTMDBID); qTMDBID != "" {
		tmdbID.Valid = true
		tmdbID.String = qTMDBID
	}

	season := model.NullInt{}
	episode := model.NullInt{}

	if qSeason := ctx.Query(torznab.ParamSeason); qSeason != "" {
		if intSeason, err := strconv.Atoi(qSeason); err == nil {
			season.Valid = true
			season.Int = intSeason

			if qEpisode := ctx.Query(torznab.ParamEpisode); qEpisode != "" {
				if intEpisode, err := strconv.Atoi(qEpisode); err == nil {
					episode.Valid = true
					episode.Int = intEpisode
				}
			}
		}
	}

	limit := model.NullUint{}
	if intLimit, limitErr := strconv.Atoi(ctx.Query(torznab.ParamLimit)); limitErr == nil && intLimit > 0 {
		limit.Valid = true
		limit.Uint = uint(intLimit)
	}

	offset := model.NullUint{}
	if intOffset, offsetErr := strconv.Atoi(ctx.Query(torznab.ParamOffset)); offsetErr == nil {
		offset.Valid = true
		offset.Uint = uint(intOffset)
	}

	result, searchErr := h.client.Search(ctx, torznab.SearchRequest{
		Profile: profile,
		Query:   ctx.Query(torznab.ParamQuery),
		Type:    tp,
		Cats:    cats,
		IMDBID:  imdbID,
		TMDBID:  tmdbID,
		Season:  season,
		Episode: episode,
		Limit:   limit,
		Offset:  offset,
	})
	if searchErr != nil {
		h.writeError(ctx, fmt.Errorf("failed to search: %w", searchErr))
		return
	}

	h.writeXML(ctx, result)
}

func (h handler) writeXML(ctx *gin.Context, obj torznab.XMLer) {
	body, err := obj.XML()
	if err != nil {
		h.writeHTTPError(ctx, fmt.Errorf("failed to encode xml: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	_, _ = ctx.Writer.Write(body)
}

func (h handler) writeError(ctx *gin.Context, err error) {
	var torznabErr torznab.Error
	if ok := errors.As(err, &torznabErr); ok {
		h.writeXML(ctx, torznabErr)
	} else {
		h.writeHTTPError(ctx, err)
	}
}

func (handler) writeHTTPError(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError

	var httpErr httpError

	if ok := errors.As(err, &httpErr); ok {
		code = httpErr.httpErrorCode()
	}

	_ = ctx.AbortWithError(code, err)
	_, _ = ctx.Writer.WriteString(err.Error() + "\n")
}

type httpError interface {
	error
	httpErrorCode() int
}

type profileNotFoundError struct {
	name string
}

func (e profileNotFoundError) Error() string {
	return fmt.Sprintf("profile not found: %s", e.name)
}

func (profileNotFoundError) httpErrorCode() int {
	return http.StatusNotFound
}

func (h handler) getProfile(c *gin.Context) (torznab.Profile, error) {
	profilePathPart := strings.ToLower(strings.Split(strings.Trim(c.Param("any"), "/"), "/")[0])
	switch profilePathPart {
	case "", "api", torznab.ProfileDefault.ID:
		return torznab.ProfileDefault, nil
	default:
		profile, ok := h.config.GetProfile(profilePathPart)
		if !ok {
			return profile, profileNotFoundError{name: profilePathPart}
		}

		return profile, nil
	}
}
