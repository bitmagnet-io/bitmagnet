package tpdb

import (
	"github.com/agnivade/levenshtein"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
)

func levenshteinCheck(target string, candidates []string, threshold uint) bool {
	normTarget := regex.NormalizeString(target)
	triedCandidates := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		normCandidate := regex.NormalizeString(candidate)
		if _, ok := triedCandidates[normCandidate]; ok {
			continue
		}
		if levenshtein.ComputeDistance(normTarget, normCandidate) <= int(threshold) {
			return true
		}
		triedCandidates[normCandidate] = struct{}{}
	}
	return false
}
