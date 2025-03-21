package httpserver

import (
  "errors"
  "fmt"
  "github.com/bitmagnet-io/bitmagnet/internal/model"
  "github.com/bitmagnet-io/bitmagnet/internal/torznab"
  "github.com/gin-gonic/gin"
  "strconv"
  "strings"
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
    caps, capsErr := h.client.Caps(ctx, profile)
    if capsErr != nil {
      h.writeError(ctx, fmt.Errorf("failed to execute caps: %w", capsErr))
      return
    }
    h.writeXML(ctx, caps)

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
  imdbId := model.NullString{}
  if qImdbId := ctx.Query(torznab.ParamImdbId); qImdbId != "" {
    imdbId.Valid = true
    imdbId.String = qImdbId
  }
  tmdbId := model.NullString{}
  if qTmdbId := ctx.Query(torznab.ParamTmdbId); qTmdbId != "" {
    tmdbId.Valid = true
    tmdbId.String = qTmdbId
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
    ImdbId:  imdbId,
    TmdbId:  tmdbId,
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

func (h handler) writeXML(ctx *gin.Context, obj torznab.Xmler) {
  body, err := obj.Xml()
  if err != nil {
    h.writeHTTPError(ctx, fmt.Errorf("failed to encode xml: %w", err))
    return
  }
  ctx.Status(200)
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

func (h handler) writeHTTPError(ctx *gin.Context, err error) {
  code := 500
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

type errProfileNotFound struct {
  name string
}

func (e errProfileNotFound) Error() string {
  return fmt.Sprintf("profile not found: %s", e.name)
}

func (e errProfileNotFound) httpErrorCode() int {
  return 404
}

func (h handler) getProfile(c *gin.Context) (torznab.Profile, error) {
  profileName := strings.ToLower(strings.Split(strings.Trim(c.Param("any"), "/"), "/")[0])
  switch profileName {
  case "", "api", torznab.ProfileDefault.Name:
    return torznab.ProfileDefault, nil
  default:
    profile, ok := h.config.GetProfile(profileName)
    if !ok {
      return profile, errProfileNotFound{name: profileName}
    }
    return profile.MergeDefaults(), nil
  }
}
