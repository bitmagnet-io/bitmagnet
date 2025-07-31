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

func Option(cfg torznab.Config, client torznab.Client) gin.OptionFunc {
	handler := Handler(cfg, client)

	return func(engine *gin.Engine) {
		engine.GET("/torznab/*any", handler)
	}
}

func Handler(cfg torznab.Config, client torznab.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		profile, err := getProfile(ctx, cfg)
		if err != nil {
			writeError(ctx, err)
			return
		}

		tp := ctx.Query(torznab.ParamType)

		switch tp {
		case "":
			writeError(ctx, torznab.Error{
				Code:        200,
				Description: fmt.Sprintf("missing parameter (%s)", torznab.ParamType),
			})

		case torznab.FunctionCaps:
			writeXML(ctx, profile.Caps())

		default:
			handleSearch(ctx, client, profile, tp)
		}
	}
}

func handleSearch(
	ctx *gin.Context,
	client torznab.Client,
	profile torznab.Profile,
	tp string,
) {
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

	result, searchErr := client.Search(ctx, torznab.SearchRequest{
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
		writeError(ctx, fmt.Errorf("failed to search: %w", searchErr))
		return
	}

	writeXML(ctx, result)
}

func writeXML(ctx *gin.Context, obj torznab.XMLer) {
	body, err := obj.XML()
	if err != nil {
		writeHTTPError(ctx, fmt.Errorf("failed to encode xml: %w", err))
		return
	}

	ctx.Status(http.StatusOK)
	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	_, _ = ctx.Writer.Write(body)
}

func writeError(ctx *gin.Context, err error) {
	var torznabErr torznab.Error
	if ok := errors.As(err, &torznabErr); ok {
		writeXML(ctx, torznabErr)
	} else {
		writeHTTPError(ctx, err)
	}
}

func writeHTTPError(ctx *gin.Context, err error) {
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

func getProfile(ctx *gin.Context, cfg torznab.Config) (torznab.Profile, error) {
	profilePathPart := strings.ToLower(strings.Split(strings.Trim(ctx.Param("any"), "/"), "/")[0])
	switch profilePathPart {
	case "", "api", torznab.ProfileDefault.ID:
		return torznab.ProfileDefault, nil
	default:
		profile, ok := cfg.GetProfile(profilePathPart)
		if !ok {
			return profile, profileNotFoundError{name: profilePathPart}
		}

		return profile, nil
	}
}
