package model

import (
	"strconv"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/keywords"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
)

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
		keywords.MustNewRexTokensFromKeywords("season", "s")...,
	).NonCaptured(),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rangeToken("sS"),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
).NonCaptured()

var episodeToken = rex.Group.Define(
	rex.Group.Composite(
		keywords.MustNewRexTokensFromKeywords("episode", "ep", "e")...,
	).NonCaptured(),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rangeToken("eE"),
).NonCaptured()

var episodesRegularTokens = rex.Group.Define(
	seasonToken,
	episodeToken.Repeat().ZeroOrOne(),
).NonCaptured()

var episodesXFormatTokens = rex.Group.Define(
	rex.Group.Define(
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)),
		rex.Chars.Runes("xX"),
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)),
	).NonCaptured(),
	rex.Group.Define(
		rex.Chars.Whitespace().Repeat().ZeroOrOne(),
		rex.Chars.Single('-'),
		rex.Chars.Whitespace().Repeat().ZeroOrOne(),
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)),
	).NonCaptured().Repeat().ZeroOrOne(),
).NonCaptured()

var EpisodesToken = rex.Group.Composite(
	episodesRegularTokens,
	episodesXFormatTokens,
)

var episodesRegex = rex.New(
	rex.Chars.Begin(),
	EpisodesToken,
	rex.Chars.End(),
).MustCompile()

func EpisodesMatchToEpisodes(match []string) Episodes {
	if len(match) < 12 {
		return nil
	}

	episodes := Episodes{}

	if match[1] != "" {
		// regular format
		seasonStart, _ := strconv.ParseInt(match[2], 10, 16)

		if match[5] == "" {
			// no episodes
			switch {
			case match[3] != "":
				// a season range
				seasonEnd, _ := strconv.ParseInt(match[3], 10, 16)
				for i := seasonStart; i <= seasonEnd; i++ {
					episodes = episodes.AddSeason(int(i))
				}
			case match[4] != "":
				// a list of seasons
				includedSeasons := strings.Split(match[1], ",")
				for _, season := range includedSeasons {
					seasonIndex, _ := strconv.ParseInt(season, 10, 16)
					episodes = episodes.AddSeason(int(seasonIndex))
				}
			default:
				// or just a single season
				episodes = episodes.AddSeason(int(seasonStart))
			}
		} else {
			// episodes
			episodeStart, _ := strconv.ParseInt(match[6], 10, 16)

			switch {
			case match[7] != "":
				// an episode range
				episodeEnd, _ := strconv.ParseInt(match[7], 10, 16)
				for i := episodeStart; i <= episodeEnd; i++ {
					episodes = episodes.AddEpisode(int(seasonStart), int(i))
				}
			case match[8] != "":
				// a list of episodes
				includedEpisodes := strings.Split(match[5], ",")
				for _, episode := range includedEpisodes {
					episodeIndex, _ := strconv.ParseInt(episode, 10, 16)
					episodes = episodes.AddEpisode(int(seasonStart), int(episodeIndex))
				}
			default:
				// a single episode
				episodes = episodes.AddEpisode(int(seasonStart), int(episodeStart))
			}
		}
	} else {
		// x format
		season, _ := strconv.ParseInt(match[9], 10, 16)
		episodeStart, _ := strconv.ParseInt(match[10], 10, 16)
		episodeEnd := episodeStart

		if match[11] != "" {
			episodeEnd, _ = strconv.ParseInt(match[11], 10, 16)
		}

		for i := episodeStart; i <= episodeEnd; i++ {
			episodes = episodes.AddEpisode(int(season), int(i))
		}
	}

	return episodes
}

func ParseEpisodes(input string) Episodes {
	return EpisodesMatchToEpisodes(episodesRegex.FindStringSubmatch(input))
}
