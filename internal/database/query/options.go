package query

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Option = func(ctx OptionBuilder) (OptionBuilder, error)

func DefaultOption() Option {
	return Options(
		Limit(10),
		WithAggregationBudget(5_000),
	)
}

func Options(o ...Option) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		for _, opt := range o {
			nextCtx, err := opt(b)
			if err != nil {
				return b, err
			}
			b = nextCtx
		}
		return b, nil
	}
}

type Criteria interface {
	Raw(c DBContext) (RawCriteria, error)
}

type TableJoinType int

const (
	TableJoinTypeInner TableJoinType = iota
	TableJoinTypeLeft
	TableJoinTypeRight
)

type TableJoin struct {
	Table schema.Tabler
	On    []field.Expr
	Type  TableJoinType
	//revive:disable-next-line:nested-structs
	Dependencies maps.InsertMap[string, struct{}]
	Required     bool
}

func Table(name string) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		return b.Table(name), nil
	}
}

func Join(fn func(*dao.Query) []TableJoin) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		joins := fn(b.Query())
		return b.Join(joins...), nil
	}
}

func RequireJoin(names ...string) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		return b.RequireJoin(names...), nil
	}
}

func SearchString(str string) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.QueryString(str), nil
	}
}

func Select(columns ...clause.Expr) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.Select(columns...), nil
	}
}

func SelectAll() Option {
	return Select(clause.Expr{SQL: "*"})
}

func Group(columns ...clause.Column) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.Group(columns...), nil
	}
}

func Limit(n uint) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.Limit(n), nil
	}
}

func Offset(n uint) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.Offset(n), nil
	}
}

func OrderBy(columns ...OrderByColumn) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.OrderBy(columns...), nil
	}
}

const QueryStringRankField = "query_string_rank"

func OrderByQueryStringRank() Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.OrderBy(OrderByColumn{
			OrderByColumn: clause.OrderByColumn{
				Column:  clause.Column{Name: QueryStringRankField},
				Desc:    true,
				Reorder: true,
			},
		}), nil
	}
}

func WithFacet(facets ...Facet) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.Facet(facets...), nil
	}
}

func Preload(fn func(query *dao.Query) []field.RelationField) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.Preload(fn(ctx.Query())...), nil
	}
}

func WithTotalCount(bl bool) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.WithTotalCount(bl), nil
	}
}

func WithHasNextPage(bl bool) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.WithHasNextPage(bl), nil
	}
}

func WithAggregationBudget(budget float64) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.WithAggregationBudget(budget), nil
	}
}

func Context(fn func(ctx context.Context) context.Context) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		return b.Context(fn), nil
	}
}

func CacheMode(mode cache.Mode) Option {
	return Context(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, cache.ModeKey, mode)
	})
}

func Cached() Option {
	return CacheMode(cache.ModeCached)
}

func CacheWarm() Option {
	return CacheMode(cache.ModeWarm)
}
