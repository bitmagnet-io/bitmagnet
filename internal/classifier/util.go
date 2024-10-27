package classifier

import (
	"github.com/agnivade/levenshtein"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
)

const levenshteinThreshold = 5

func levenshteinCheck(target string, candidates []string, threshold uint) (bool, int) {
	normTarget := regex.NormalizeString(target)
	triedCandidates := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		normCandidate := regex.NormalizeString(candidate)
		if _, ok := triedCandidates[normCandidate]; ok {
			continue
		}
		distance := levenshtein.ComputeDistance(normTarget, normCandidate)
		if distance <= int(threshold) {
			return true, distance
		}
		triedCandidates[normCandidate] = struct{}{}
	}
	return false, -1
}
