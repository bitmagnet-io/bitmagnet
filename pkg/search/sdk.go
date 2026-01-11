package search

import "github.com/bitmagnet-io/bitmagnet/internal/search"

var (
	ParseCriteria = search.ParseCriteria
)

type (
	Criteria            = search.Criteria
	CriteriaContentRef  = search.CriteriaContentRef
	CriteriaContentType = search.CriteriaContentType
	CriteriaGenre       = search.CriteriaGenre
	CriteriaInfoHash    = search.CriteriaInfoHash
	CriteriaLanguage    = search.CriteriaLanguage
	And                 = search.And
	Or                  = search.Or
	Not                 = search.Not
)
