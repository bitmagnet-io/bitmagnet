package model

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/keywords"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
)

func rangeToken(runes string, startName, dashEndName, commaEndName string) dialect.Token {
	return rex.Group.Define(
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName(startName),
		rex.Group.Composite(
			rex.Group.Define(
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Chars.Single('-'),
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Group.NonCaptured(
					rex.Chars.Runes(runes).Repeat().ZeroOrOne(),
					rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				).Repeat().ZeroOrOne(),
				rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName(dashEndName),
			).NonCaptured(),
			rex.Group.Define(
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Chars.Single(','),
				rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				rex.Group.NonCaptured(
					rex.Chars.Runes(runes).Repeat().ZeroOrOne(),
					rex.Chars.Whitespace().Repeat().ZeroOrOne(),
				).Repeat().ZeroOrOne(),
				rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName(commaEndName),
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
	rangeToken("sS", "seasonStart", "seasonDashEnd", "seasonCommaEnd"),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
).NonCaptured()

var episodeToken = rex.Group.Define(
	rex.Group.Composite(
		keywords.MustNewRexTokensFromKeywords("episode", "ep", "e")...,
	).NonCaptured(),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rangeToken("eE", "episodeStart", "episodeDashEnd", "episodeCommaEnd"),
).NonCaptured()

var episodeDashToken = rex.Group.Define(
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rex.Chars.Single('-'),
	rex.Chars.Whitespace().Repeat().ZeroOrOne(),
	rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName("episodeStart"),
).NonCaptured()

var episodesRegularTokens = rex.Group.Define(
	seasonToken,
	rex.Group.Composite(
		episodeToken,
		episodeDashToken,
	).NonCaptured().Repeat().ZeroOrOne(),
).NonCaptured()

var episodesXFormatTokens = rex.Group.Define(
	rex.Group.Define(
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName("xSeason"),
		rex.Chars.Runes("xX"),
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName("xEpisodeStart"),
	).NonCaptured(),
	rex.Group.Define(
		rex.Chars.Whitespace().Repeat().ZeroOrOne(),
		rex.Chars.Single('-'),
		rex.Chars.Whitespace().Repeat().ZeroOrOne(),
		rex.Group.Define(rex.Chars.Digits().Repeat().Between(1, 2)).WithName("xEpisodeEnd"),
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

func namedMatch(re *regexp.Regexp, match []string, name string) string {
	for i, n := range re.SubexpNames() {
		if n == name && i < len(match) {
			return match[i]
		}
	}
	return ""
}

// EpisodesMatchToEpisodes converts a regex submatch to Episodes using named groups.
// The re parameter should be a regex compiled from EpisodesToken.
func EpisodesMatchToEpisodes(re *regexp.Regexp, match []string) Episodes {
	if match == nil {
		return nil
	}

	nm := func(name string) string {
		return namedMatch(re, match, name)
	}

	episodes := Episodes{}

	if nm("seasonStart") != "" {
		// regular format
		seasonStart, _ := strconv.ParseInt(nm("seasonStart"), 10, 16)

		if nm("episodeStart") == "" {
			// no episodes
			switch {
			case nm("seasonDashEnd") != "":
				// a season range
				seasonEnd, _ := strconv.ParseInt(nm("seasonDashEnd"), 10, 16)
				for i := seasonStart; i <= seasonEnd; i++ {
					episodes = episodes.AddSeason(int(i))
				}
			case nm("seasonCommaEnd") != "":
				// a list of seasons - find the parent range group containing commas
				var seasonRange string
				for i, n := range re.SubexpNames() {
					if n == "seasonStart" && i < len(match) {
						for j := i - 1; j >= 0; j-- {
							if match[j] != "" && strings.Contains(match[j], ",") {
								seasonRange = match[j]
								break
							}
						}
						break
					}
				}
				includedSeasons := strings.Split(seasonRange, ",")
				for _, season := range includedSeasons {
					season = strings.TrimSpace(season)
					season = strings.TrimLeft(season, "sS ")
					seasonIndex, _ := strconv.ParseInt(season, 10, 16)
					episodes = episodes.AddSeason(int(seasonIndex))
				}
			default:
				// or just a single season
				episodes = episodes.AddSeason(int(seasonStart))
			}
		} else {
			// episodes
			episodeStart, _ := strconv.ParseInt(nm("episodeStart"), 10, 16)

			switch {
			case nm("episodeDashEnd") != "":
				// an episode range
				episodeEnd, _ := strconv.ParseInt(nm("episodeDashEnd"), 10, 16)
				for i := episodeStart; i <= episodeEnd; i++ {
					episodes = episodes.AddEpisode(int(seasonStart), int(i))
				}
			case nm("episodeCommaEnd") != "":
				// a list of episodes - find the parent range group containing commas
				var episodeRange string
				for i, n := range re.SubexpNames() {
					if n == "episodeStart" && i < len(match) {
						for j := i - 1; j >= 0; j-- {
							if match[j] != "" && strings.Contains(match[j], ",") {
								episodeRange = match[j]
								break
							}
						}
						break
					}
				}
				includedEpisodes := strings.Split(episodeRange, ",")
				for _, episode := range includedEpisodes {
					episode = strings.TrimSpace(episode)
					episode = strings.TrimLeft(episode, "eE ")
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
		season, _ := strconv.ParseInt(nm("xSeason"), 10, 16)
		episodeStart, _ := strconv.ParseInt(nm("xEpisodeStart"), 10, 16)
		episodeEnd := episodeStart

		if nm("xEpisodeEnd") != "" {
			episodeEnd, _ = strconv.ParseInt(nm("xEpisodeEnd"), 10, 16)
		}

		for i := episodeStart; i <= episodeEnd; i++ {
			episodes = episodes.AddEpisode(int(season), int(i))
		}
	}

	return episodes
}

func ParseEpisodes(input string) Episodes {
	return EpisodesMatchToEpisodes(episodesRegex, episodesRegex.FindStringSubmatch(input))
}
