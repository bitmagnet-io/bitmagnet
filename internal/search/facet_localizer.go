package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

func NewFacetLocalizer(ref ref.Ref, i18n *i18n.Bundle) FacetLocalizer {
	return FacetLocalizer{
		ref:  ref,
		i18n: i18n,
	}
}

type FacetLocalizer struct {
	ref  ref.Ref
	i18n *i18n.Bundle
}

func (f FacetLocalizer) Label(facet Facet, acceptLanguage []string) string {
	localizer := i18n.NewLocalizer(f.i18n, acceptLanguage...)
	result, _ := localizer.LocalizeMessage(&i18n.Message{
		ID: f.ref.String() + "." + facet.String(),
	})
	if result == "" {
		return facet.Label()
	}
	return result
}
