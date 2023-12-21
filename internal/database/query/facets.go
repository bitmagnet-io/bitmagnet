package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"sort"
	"sync"
)

type FacetConfig interface {
	Key() string
	Label() string
	Logic() model.FacetLogic
	Filter() FacetFilter
	IsAggregated() bool
	AggregationOption(b OptionBuilder) (OptionBuilder, error)
}

type Facet interface {
	FacetConfig
	Aggregate(ctx FacetContext) (AggregationItems, error)
	Criteria() []Criteria
}

type FacetFilter map[string]struct{}

// Values allows iteration over deterministically sorted filter values, which helps with query caching.
func (f FacetFilter) Values() []string {
	values := make([]string, 0, len(f))
	for v := range f {
		values = append(values, v)
	}
	sort.Strings(values)
	return values
}

type facetConfig struct {
	key                string
	label              string
	logic              model.FacetLogic
	filter             FacetFilter
	aggregate          bool
	aggregationOptions []Option
}

type FacetOption func(facetConfig) facetConfig

func NewFacetConfig(options ...FacetOption) FacetConfig {
	c := facetConfig{}
	for _, option := range options {
		c = option(c)
	}
	return c
}

type AggregationItem struct {
	Label string
	Count uint
}

type AggregationItems = map[string]AggregationItem

type AggregationGroup struct {
	Label string
	Logic model.FacetLogic
	Items AggregationItems
}

type Aggregations = maps.StringMap[AggregationGroup]

type FacetContext interface {
	DbContext
	Context() context.Context
	NewAggregationQuery(options ...Option) (SubQuery, error)
}

type facetContext struct {
	optionBuilder OptionBuilder
	ctx           context.Context
}

func (ctx facetContext) Query() *dao.Query {
	return ctx.optionBuilder.Query()
}

func (ctx facetContext) TableName() string {
	return ctx.optionBuilder.TableName()
}

func (ctx facetContext) NewSubQuery(c context.Context) SubQuery {
	return ctx.optionBuilder.NewSubQuery(c)
}

func (ctx facetContext) Context() context.Context {
	return ctx.ctx
}

func (ctx facetContext) NewAggregationQuery(options ...Option) (SubQuery, error) {
	subCtx, subErr := Options(options...)(ctx.optionBuilder)
	if subErr != nil {
		return nil, subErr
	}
	sq := ctx.optionBuilder.NewSubQuery(ctx.Context())
	applyErr := subCtx.applyPre(sq)
	if applyErr != nil {
		return nil, applyErr
	}
	return sq, nil
}

func FacetHasKey(key string) FacetOption {
	return func(a facetConfig) facetConfig {
		a.key = key
		return a
	}
}

func FacetHasLabel(label string) FacetOption {
	return func(a facetConfig) facetConfig {
		a.label = label
		return a
	}
}

func FacetUsesLogic(logic model.FacetLogic) FacetOption {
	return func(a facetConfig) facetConfig {
		a.logic = logic
		return a
	}
}

func FacetUsesAndLogic() FacetOption {
	return FacetUsesLogic(model.FacetLogicAnd)
}

func FacetUsesOrLogic() FacetOption {
	return FacetUsesLogic(model.FacetLogicOr)
}

func FacetHasFilter(filter FacetFilter) FacetOption {
	return func(c facetConfig) facetConfig {
		c.filter = filter
		return c
	}
}

func FacetIsAggregated() FacetOption {
	return func(c facetConfig) facetConfig {
		c.aggregate = true
		return c
	}
}

func FacetHasAggregationOption(options ...Option) FacetOption {
	return func(c facetConfig) facetConfig {
		c.aggregationOptions = append(c.aggregationOptions, options...)
		return c
	}
}

func (c facetConfig) Key() string {
	return c.key
}

func (c facetConfig) Label() string {
	return c.label
}

func (c facetConfig) Logic() model.FacetLogic {
	return c.logic
}

func (c facetConfig) IsAggregated() bool {
	return c.aggregate
}

func (c facetConfig) AggregationOption(b OptionBuilder) (OptionBuilder, error) {
	return Options(c.aggregationOptions...)(b)
}

func (c facetConfig) Filter() FacetFilter {
	return c.filter
}

func (b optionBuilder) createFacetsFilterCriteria() (c Criteria, err error) {
	cs := make([]Criteria, 0, len(b.facets))
	for _, facet := range b.facets {
		cr := facet.Criteria()
		switch facet.Logic() {
		case model.FacetLogicAnd:
			cs = append(cs, AndCriteria{cr})
		case model.FacetLogicOr:
			if b.currentFacet != facet.Key() {
				cs = append(cs, OrCriteria{cr})
			}
		}
	}
	return AndCriteria{cs}, nil
}

// when aggregating with or logic the current facet's filter should be ignored
func withCurrentFacet(facetKey string) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		return b.withCurrentFacet(facetKey), nil
	}
}

func (b optionBuilder) calculateAggregations(ctx context.Context) (Aggregations, error) {
	aggregations := make(Aggregations, len(b.facets))
	wg := sync.WaitGroup{}
	wg.Add(len(b.facets))
	mtx := sync.Mutex{}
	var errs []error
	for _, facet := range b.facets {
		go (func(facet Facet) {
			defer wg.Done()
			if !facet.IsAggregated() {
				return
			}
			aggBuilder, aggBuilderErr := Options(facet.AggregationOption, withCurrentFacet(facet.Key()))(b)
			if aggBuilderErr != nil {
				errs = append(errs, fmt.Errorf("failed to create aggregation option for key '%s': %w", facet.Key(), aggBuilderErr))
				return
			}
			aggCtx := facetContext{
				optionBuilder: aggBuilder,
				ctx:           ctx,
			}
			aggregation, aggregateErr := facet.Aggregate(aggCtx)
			mtx.Lock()
			defer mtx.Unlock()
			if aggregateErr != nil {
				errs = append(errs, fmt.Errorf("failed to aggregate key '%s': %w", facet.Key(), aggregateErr))
			} else {
				aggregations[facet.Key()] = AggregationGroup{
					Label: facet.Label(),
					Logic: facet.Logic(),
					Items: aggregation,
				}
			}
		})(facet)
	}
	wg.Wait()
	return aggregations, errors.Join(errs...)
}
