package query

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type Option = func(ctx OptionBuilder) (OptionBuilder, error)

func DefaultOption() Option {
	return Options(
		Limit(10),
	)
}

func Options(o ...Option) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		for _, opt := range o {
			if nextCtx, err := opt(b); err != nil {
				return b, err
			} else {
				b = nextCtx
			}
		}
		return b, nil
	}
}

type Criteria interface {
	Raw(c DbContext) (RawCriteria, error)
}

type TableJoinType int

const (
	TableJoinTypeInner TableJoinType = iota
	TableJoinTypeLeft
	TableJoinTypeRight
)

type TableJoin struct {
	Table        schema.Tabler
	On           []field.Expr
	Type         TableJoinType
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

func QueryString(str string) Option {
	query := fts.AppQueryToTsquery(str)
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		if len(query) == 0 {
			return ctx.Select(clause.Expr{
				SQL: "0 AS " + queryStringRankField,
			}), nil
		}
		c, err := GenCriteria(func(ctx DbContext) (Criteria, error) {
			return DbCriteria{
				Sql: strings.Join([]string{
					ctx.TableName() + ".tsv @@ ?::tsquery",
				}, " OR "),
				Args: []interface{}{
					query,
				},
			}, nil
		}).Raw(ctx)
		if err != nil {
			return ctx, err
		}
		ctx = ctx.Scope(func(dao SubQuery) error {
			dao.UnderlyingDB().Where(c.Query, c.Args...)
			return nil
		}).RequireJoin(ctx.TableName()).Select(clause.Expr{
			SQL: "ts_rank_cd(" + ctx.TableName() + ".tsv, ?::tsquery) AS " + queryStringRankField,
			Vars: []interface{}{
				query,
			},
		})
		return ctx, nil
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

func OrderBy(columns ...clause.OrderByColumn) Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.OrderBy(columns...), nil
	}
}

const queryStringRankField = "query_string_rank"

func OrderByQueryStringRank() Option {
	return func(ctx OptionBuilder) (OptionBuilder, error) {
		return ctx.OrderBy(clause.OrderByColumn{
			Column:  clause.Column{Name: queryStringRankField},
			Desc:    true,
			Reorder: true,
		}), nil
	}
}

func OrderByColumn(field string, desc bool) Option {
	return OrderBy(clause.OrderByColumn{
		Column: clause.Column{Name: field},
		Desc:   desc,
	})
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
