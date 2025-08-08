package metrics

import (
	"net/url"

	"github.com/hashicorp/go-metrics"
)

// Label represents a grouping dimension for a metric.
type Label string

// Value creates a LabelValue by combining this Label with the provided value string.
func (l Label) Value(value string) LabelValue {
	return LabelValue{
		Label: l,
		Value: value,
	}
}

// LabelValue represents a value for a Label.
type LabelValue struct {
	Label
	Value string
}

func (lv LabelValue) encodeValue() string {
	return url.QueryEscape(lv.Value)
}

// encode encodes a label value for use by the go-metrics library.
// URL query escaping is used to allow semicolons and equals to act as delimiters.
func (lv LabelValue) encode() string {
	return lv.Label.encode() + "=" + lv.encodeValue()
}

func (lv LabelValue) metricsLabel() metrics.Label {
	return metrics.Label{
		Name:  lv.Label.encode(),
		Value: lv.encodeValue(),
	}
}

func (l Label) encode() string {
	return url.QueryEscape(string(l))
}
