package metrics

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/go-metrics"
)

var spaceReplacer = strings.NewReplacer(" ", "_")

// flattenRefLabels flattens the key for formatting along with its labelValues, removes spaces.
func flattenRefLabels(parts []string, labels ...metrics.Label) string {
	buf := &bytes.Buffer{}

	joined := strings.Join(parts[1:], ".")

	_, _ = spaceReplacer.WriteString(buf, joined)

	sortLabels(labels)

	for _, label := range labels {
		_, _ = spaceReplacer.WriteString(buf, fmt.Sprintf(";%s=%s", label.Name, label.Value))
	}

	return buf.String()
}

// checkHash checks if a metric hash is a match for the specified key and all label values.
func checkHash(hash string, typeRef Ref, labels []LabelValue) bool {
	hashParts := strings.SplitN(hash, ";", 2)

	if hashParts[0] != typeRef.String() {
		return false
	}

	var hashLabelParts []string
	if len(hashParts) > 1 {
		hashLabelParts = strings.Split(hashParts[1], ";")
	}

	hashLabelMap := make(map[string]struct{}, len(hashLabelParts))

	for _, labelPart := range hashLabelParts {
		hashLabelMap[labelPart] = struct{}{}
	}

	// for Gauge metrics, the labels must be an exact match
	var labelMap map[string]struct{}

	if typeRef.Type == TypeGauge {
		labelMap = make(map[string]struct{}, len(labels))
	}

	for _, label := range labels {
		encoded := label.encode()
		if _, ok := hashLabelMap[encoded]; !ok {
			return false
		}

		if labelMap != nil {
			labelMap[encoded] = struct{}{}
		}
	}

	if labelMap != nil && len(labelMap) != len(hashLabelMap) {
		return false
	}

	return true
}

// metricsLabels converts a slice of LabelValue to a slice of metrics.Label.
func metricsLabels(labelValues []LabelValue) []metrics.Label {
	result := make([]metrics.Label, len(labelValues))

	for i, labelValue := range labelValues {
		result[i] = labelValue.metricsLabel()
	}

	sortLabels(result)

	return result
}

func sortLabels(labels []metrics.Label) {
	slices.SortFunc(labels, func(a, b metrics.Label) int {
		result := cmp.Compare(a.Name, b.Name)
		if result == 0 {
			result = cmp.Compare(a.Value, b.Value)
		}

		return result
	})
}
