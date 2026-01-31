package metrics

import (
	"net/url"
	"strings"
)

// Averages is a map of string keys to Average values.
type Averages map[string]Average

// Average returns an Average for the provided Key, optionally filtered by LabelValue.
func (a Averages) Average(ref Ref, labels ...LabelValue) Average {
	var result Average

	for hash, avg := range a {
		if checkHash(hash, ref, labels) {
			result.Count += avg.Count
			result.Sum += avg.Sum
		}
	}

	return result
}

// AveragesByLabel returns all averages for the provided key, grouped by the value of the provided Label.
func (a Averages) AveragesByLabel(ref Ref, label Label) map[string]Average {
	strRef := ref.String()
	result := make(map[string]Average)

	for hash, avg := range a {
		hashParts := strings.Split(hash, ";")
		if len(hashParts) < 2 || hashParts[0] != strRef {
			continue
		}

		for _, labelPart := range hashParts[1:] {
			labelSplit := strings.SplitN(labelPart, "=", 2)
			if len(labelSplit) < 2 || labelSplit[0] != label.encode() {
				continue
			}

			unescaped, err := url.QueryUnescape(labelSplit[1])
			if err != nil {
				continue
			}

			current := result[unescaped]
			current.Count += avg.Count
			current.Sum += avg.Sum
			result[unescaped] = current
		}
	}

	return result
}

// Average stores an accumulation of values and their count for calculating means
type Average struct {
	Count int
	Sum   float64
}

// Value returns the average value.
func (a Average) Value() float64 {
	if a.Count == 0 {
		return 0
	}

	return a.Sum / float64(a.Count)
}

func (t *sink) Averages() Averages {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.averagesUnlocked()
}

func (t *sink) averagesUnlocked() Averages {
	averages := make(Averages, len(t.averages))

	for k, v := range t.averages {
		averages[k] = v
	}

	return averages
}
