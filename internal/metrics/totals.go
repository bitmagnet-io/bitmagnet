package metrics

import (
	"net/url"
	"strings"
)

// Totals is a map of metric hashes to their total value.
type Totals map[string]int

// Value returns the total value for the provided Key, optionally filtered by LabelValue.
func (t Totals) Value(typeRef Ref, labels ...LabelValue) int {
	total := 0

	for hash, count := range t {
		if checkHash(hash, typeRef, labels) {
			total += count
		}
	}

	return total
}

// ValuesByLabel returns all totals for the provided key, grouped by the value of the provided Label.
func (t Totals) ValuesByLabel(typeRef Ref, label Label) map[string]int {
	strRef := typeRef.String()
	result := make(map[string]int)

	for hash, count := range t {
		hashParts := strings.Split(hash, ";")
		if len(hashParts) < 2 || hashParts[0] != strRef || (typeRef.Type == TypeGauge && len(hashParts) != 2) {
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

			result[unescaped] += count
		}
	}

	return result
}

func (t *sink) Totals() Totals {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.totalsUnlocked()
}

func (t *sink) totalsUnlocked() Totals {
	totals := make(Totals, len(t.totals))

	for k, v := range t.totals {
		totals[k] = v
	}

	return totals
}
