package classifier

import (
	"encoding/json"
	"regexp"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"go.uber.org/zap"
)

type dependencies struct {
	search     LocalSearch
	tmdbClient tmdb.Client
	_logger    *zap.SugaredLogger
	logger		 *zap.SugaredLogger
}

func (d *dependencies) CleanObj(o interface{}) map[string]any {
	var isEmptyString = regexp.MustCompile("^(?:[0\\s]*|0001-01-01T00:00:00Z)$")
	var m map[string]any
	jhint,_ := json.Marshal(o)
	json.Unmarshal(jhint, &m)
	for k,v := range m {
		if s, sOk := v.(string); sOk && isEmptyString.MatchString(s) {
			delete(m, k)
			continue
		}		
		if a, aOk := v.([]any); aOk && len(a) == 0 {
			delete(m, k)
			continue
		}
		if f, fOk := v.(float64); fOk && f == 0.0 {
			delete(m, k)
			continue
		}
		d, dOk := v.(model.Date)		
		if dOk && (d.Year == 0 || d.Month == 0 || d.Day == 0) {
			delete(m, k)
			continue
		}
		if v == nil {
			delete(m, k)
			continue
		}
	}
	return m
}