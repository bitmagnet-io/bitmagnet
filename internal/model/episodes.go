package model

import (
	"fmt"
	"sort"
	"strings"
)

type Episodes map[int]map[int]struct{}

func (e Episodes) String() string {
	seasonEntries := e.SeasonEntries()

	var wholeSeasons []int

	for _, season := range seasonEntries {
		if len(season.Episodes) == 0 {
			wholeSeasons = append(wholeSeasons, season.Season)
		}
	}

	contiguousWholeSeasons := getContiguousRanges(wholeSeasons)
	partsMap := make(map[int]string)

	for _, r := range contiguousWholeSeasons {
		partsMap[r.start] = "S" + r.String()
	}

	for _, entry := range seasonEntries {
		if len(entry.Episodes) > 0 {
			partsMap[entry.Season] = entry.String()
		}
	}

	parts := make([]string, 0, len(partsMap))

	for _, entry := range seasonEntries {
		if part, ok := partsMap[entry.Season]; ok {
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, ", ")
}

func (e Episodes) HasEpisode(season, episode int) bool {
	epMap, epMapOk := e[season]
	if !epMapOk {
		return false
	}

	if len(epMap) == 0 {
		return true
	}

	_, epOk := epMap[episode]

	return epOk
}

func (e Episodes) AddEpisode(season, episode int) Episodes {
	_, epMapOk := e[season]
	if !epMapOk {
		e[season] = make(map[int]struct{})
	}

	e[season][episode] = struct{}{}

	return e
}

func (e Episodes) AddSeason(season int) Episodes {
	e[season] = make(map[int]struct{})
	return e
}

func (e Episodes) SeasonEntries() []Season {
	seasonEpisodes := make([]Season, 0, len(e))

	for season, episodes := range e {
		se := Season{Season: season}

		if len(episodes) > 0 {
			eps := make([]int, 0, len(episodes))
			for episode := range episodes {
				eps = append(eps, episode)
			}

			sort.Slice(eps, func(i, j int) bool {
				return eps[i] < eps[j]
			})

			se.Episodes = eps
		}

		seasonEpisodes = append(seasonEpisodes, se)
	}

	sort.Slice(seasonEpisodes, func(i, j int) bool {
		return seasonEpisodes[i].Season < seasonEpisodes[j].Season
	})

	return seasonEpisodes
}

type Season struct {
	Season   int
	Episodes []int `json:"omitempty"`
}

func (se Season) String() string {
	str := fmt.Sprintf("S%02d", se.Season)

	if len(se.Episodes) > 0 {
		contiguousRanges := getContiguousRanges(se.Episodes)
		epParts := make([]string, 0, len(contiguousRanges))

		for _, r := range contiguousRanges {
			epParts = append(epParts, "E"+r.String())
		}

		str += strings.Join(epParts, ",")
	}

	return str
}

func getContiguousRanges(orderedInts []int) []contiguousRange {
	//nolint:prealloc
	var contiguousRanges []contiguousRange

	var currentRange contiguousRange

	for i, v := range orderedInts {
		if i == 0 {
			currentRange = contiguousRange{start: v, end: v}
			continue
		}

		if v == currentRange.end+1 {
			currentRange.end = v
			continue
		}

		contiguousRanges = append(contiguousRanges, currentRange)
		currentRange = contiguousRange{start: v, end: v}
	}

	contiguousRanges = append(contiguousRanges, currentRange)

	return contiguousRanges
}

type contiguousRange struct {
	start int
	end   int
}

func (r contiguousRange) String() string {
	if r.start == r.end {
		return fmt.Sprintf("%02d", r.start)
	}

	return fmt.Sprintf("%02d-%02d", r.start, r.end)
}
