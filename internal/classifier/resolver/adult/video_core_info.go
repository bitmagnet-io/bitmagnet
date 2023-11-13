package adult

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
	"strconv"
	"strings"
)

type CoreInfo struct {
	Title    string
	Year     model.Year
	Episodes model.Episodes
}

var titleTokens = []dialect.Token{
	rex.Group.Define(
		rex.Group.Composite(
			rex.Group.NonCaptured(
				regex.AnyWordChar().Repeat().OneOrMore(),
				rex.Group.NonCaptured(
					rex.Chars.Single('-'), regex.AnyWordChar().Repeat().OneOrMore(),
				).Repeat().ZeroOrMore(),
			),
			regex.AnyNonWordChar().Repeat().OneOrMore(),
		).NonCaptured().Repeat().OneOrMore(),
		rex.Group.Composite(
			regex.AnyNonWordChar().Repeat().OneOrMore(),
			rex.Chars.End(),
		).NonCaptured(),
	),
}

var titleRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(titleTokens...),
).MustCompile()

var yearTokens = []dialect.Token{
	rex.Group.NonCaptured(rex.Common.NotClass(rex.Chars.WordCharacter()).Repeat().ZeroOrMore()),
	rex.Group.Define(
		rex.Group.Composite(
			rex.Common.Text("18"), rex.Common.Text("19"), rex.Common.Text("20"),
		).NonCaptured(),
		rex.Chars.Digits().Repeat().Exactly(2),
	),
	rex.Group.Composite(
		rex.Common.NotClass(rex.Chars.WordCharacter()),
		rex.Chars.End(),
	).NonCaptured(),
}

var titleYearRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(rex.Group.NonCaptured(titleTokens...), rex.Group.NonCaptured(yearTokens...)),
).MustCompile()

func rangeToken(runes string) dialect.Token {
	return rex.Group.Define(
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)),
		rex.Group.Composite(
			rex.Group.Define(
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Chars.Single('-'),
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Group.NonCaptured(
					rex.Chars.Runes(runes).Repeat().ZeroOrOne(),
					rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				).Repeat().ZeroOrOne(),
				rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)),
			).NonCaptured(),
			rex.Group.Define(
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Chars.Single(','),
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Group.NonCaptured(
					rex.Chars.Runes(runes).Repeat().ZeroOrOne(),
					rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				).Repeat().ZeroOrOne(),
				rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)),
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
			).NonCaptured().Repeat().OneOrMore(),
		).NonCaptured().Repeat().ZeroOrOne(),
	)
}

var seasonToken = rex.Group.Define(
	rex.Group.Composite(
		regex.RegexTokensFromNames("season", "s")...,
	).NonCaptured(),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rangeToken("sS"),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
).NonCaptured()

var episodeToken = rex.Group.Define(
	rex.Group.Composite(
		regex.RegexTokensFromNames("episode", "ep", "e")...,
	).NonCaptured(),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rangeToken("eE"),
).NonCaptured()

var episodesTokens = rex.Group.Define(
	seasonToken,
	episodeToken.Repeat().ZeroOrOne(),
).NonCaptured()

var titleEpisodesRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(rex.Group.NonCaptured(titleTokens...), episodesTokens),
).MustCompile()

var separatorToken = rex.Chars.Runes(" ._")

var titlePartRegex = rex.New(
	separatorToken.Repeat().ZeroOrOne(),
	rex.Group.Define(regex.WordToken()),
	separatorToken.Repeat().ZeroOrOne(),
).MustCompile()

var trimTitleRegex = rex.New(
	rex.Chars.Begin(),
	rex.Group.NonCaptured(
		rex.Chars.Single('['),
		rex.Common.NotClass(rex.Chars.Single(']')).Repeat().OneOrMore(),
		rex.Chars.Single(']'),
	).Repeat().ZeroOrOne(),
	regex.AnyNonWordChar().Repeat().ZeroOrMore(),
	rex.Group.Define(
		regex.WordToken(),
		rex.Group.NonCaptured(
			rex.Chars.Any(),
			regex.WordToken(),
		).Repeat().ZeroOrMore(),
	),
	regex.AnyNonWordChar().Repeat().ZeroOrMore(),
	rex.Chars.End(),
).MustCompile()

func cleanTitle(title string) string {
	title = titlePartRegex.ReplaceAllStringFunc(title, func(s string) string {
		partMatch := titlePartRegex.FindStringSubmatch(s)
		if partMatch == nil {
			return ""
		}
		return partMatch[1] + " "
	})
	title = trimTitleRegex.ReplaceAllString(title, "$1")
	return title
}

func parseTitleYear(input string) (CoreInfo, string, error) {
	if match := titleYearRegex.FindStringSubmatch(input); match != nil {
		yearMatch, _ := strconv.ParseUint(match[2], 10, 16)
		title := cleanTitle(match[1])
		if title != "" {
			return CoreInfo{
				Title: title,
				Year:  model.Year(yearMatch),
			}, input[len(match[0]):], nil
		}
	}
	return CoreInfo{}, "", resolver.ErrNoMatch
}

func parseTitle(input string) (CoreInfo, string, error) {
	if match := titleRegex.FindStringSubmatch(input); match != nil {
		title := cleanTitle(match[1])
		if title != "" {
			return CoreInfo{
				Title: title,
			}, input[len(match[0]):], nil
		}
	}
	return CoreInfo{}, "", resolver.ErrNoMatch
}

func parseTitleYearEpisodes(input string) (CoreInfo, string, error) {
	if match := titleEpisodesRegex.FindStringSubmatch(input); match != nil {
		title := match[1]
		year := model.Year(0)
		if i, _, err := parseTitleYear(title); err == nil {
			title = i.Title
			year = i.Year
		} else {
			title = cleanTitle(title)
		}
		episodes := model.Episodes{}
		seasonStart, _ := strconv.ParseInt(match[3], 10, 16)
		if match[6] == "" {
			// no episodes
			if match[4] != "" {
				// a season range
				seasonEnd, _ := strconv.ParseInt(match[4], 10, 16)
				for i := seasonStart; i <= seasonEnd; i++ {
					episodes = episodes.AddSeason(int(i))
				}
			} else if match[5] != "" {
				// a list of seasons
				includedSeasons := strings.Split(match[2], ",")
				for _, season := range includedSeasons {
					seasonIndex, _ := strconv.ParseInt(season, 10, 16)
					episodes = episodes.AddSeason(int(seasonIndex))
				}
			} else {
				// or just a single season
				episodes = episodes.AddSeason(int(seasonStart))
			}
		} else {
			// episodes
			episodeStart, _ := strconv.ParseInt(match[7], 10, 16)
			if match[8] != "" {
				// an episode range
				episodeEnd, _ := strconv.ParseInt(match[8], 10, 16)
				for i := episodeStart; i <= episodeEnd; i++ {
					episodes = episodes.AddEpisode(int(seasonStart), int(i))
				}
			} else if match[9] != "" {
				// a list of episodes
				includedEpisodes := strings.Split(match[6], ",")
				for _, episode := range includedEpisodes {
					episodeIndex, _ := strconv.ParseInt(episode, 10, 16)
					episodes = episodes.AddEpisode(int(seasonStart), int(episodeIndex))
				}
			} else {
				// a single episode
				episodes = episodes.AddEpisode(int(seasonStart), int(episodeStart))
			}
		}
		return CoreInfo{
			Title:    title,
			Year:     year,
			Episodes: episodes,
		}, input[len(match[0]):], nil
	}
	return CoreInfo{}, "", resolver.ErrNoMatch
}

func ParseVideoCoreInfo(contentType model.NullContentType, input string) (CoreInfo, string, error) {
	if !contentType.Valid || contentType.ContentType == model.ContentTypeTvShow {
		if info, rest, err := parseTitleYearEpisodes(input); err == nil {
			return info, rest, nil
		}
	}
	if info, rest, err := parseTitleYear(input); err == nil {
		return info, rest, nil
	}
	if info, rest, err := parseTitle(input); err == nil {
		return info, rest, nil
	}
	return CoreInfo{}, "", resolver.ErrNoMatch
}
