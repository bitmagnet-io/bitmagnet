package query

import (
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"strings"
)

func Where(conditions ...Criteria) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		rawCriteria := make([]RawCriteria, 0, len(conditions))
		joins := maps.NewInsertMap[string, struct{}]()
		for _, c := range conditions {
			rc, rawCriteriaErr := c.Raw(b)
			if rawCriteriaErr != nil {
				return b, rawCriteriaErr
			}
			rawCriteria = append(rawCriteria, rc)
			joins.SetEntries(rc.Joins.Entries()...)
		}
		b = b.Scope(func(db *gorm.DB) error {
			for _, raw := range rawCriteria {
				db.Where(raw.Query, raw.Args...)
			}
			return nil
		})
		b = b.RequireJoin(joins.Keys()...)
		return b, nil
	}
}

func And(criteria ...Criteria) Criteria {
	return AndCriteria{
		Criteria: criteria,
	}
}

func Or(criteria ...Criteria) Criteria {
	return OrCriteria{
		Criteria: criteria,
	}
}

func Not(criteria ...Criteria) Criteria {
	return NotCriteria{
		Criteria: criteria,
	}
}

type RawCriteria struct {
	Query interface{}
	Args  []interface{}
	Joins maps.InsertMap[string, struct{}]
}

func (c RawCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	return c, nil
}

type DaoCriteria struct {
	Conditions func(ctx DbContext) ([]field.Expr, error)
	Joins      maps.InsertMap[string, struct{}]
}

func (c DaoCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	// todo Don't reference model
	sq := ctx.Query().Torrent.UnderlyingDB()
	conditions, conditionsErr := c.Conditions(ctx)
	if conditionsErr != nil {
		return RawCriteria{}, conditionsErr
	}
	for _, condition := range conditions {
		sq = sq.Where(condition.RawExpr())
	}
	return RawCriteria{
		Query: sq,
		Joins: c.Joins,
	}, nil
}

type OrCriteria struct {
	Criteria []Criteria
}

func (c OrCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	joins := maps.NewInsertMap[string, struct{}]()
	sq := ctx.Query().Torrent.UnderlyingDB()
	for _, c := range c.Criteria {
		rc, rawCriteriaErr := c.Raw(ctx)
		if rawCriteriaErr != nil {
			return RawCriteria{}, rawCriteriaErr
		}
		joins.SetEntries(rc.Joins.Entries()...)
		sq = sq.Or(rc.Query, rc.Args...)
	}
	return RawCriteria{
		Joins: joins,
		Query: sq,
	}, nil
}

type AndCriteria struct {
	Criteria []Criteria
}

func (c AndCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	joins := maps.NewInsertMap[string, struct{}]()
	sq := ctx.Query().Torrent.UnderlyingDB()
	for _, c := range c.Criteria {
		rc, rawCriteriaErr := c.Raw(ctx)
		if rawCriteriaErr != nil {
			return RawCriteria{}, rawCriteriaErr
		}
		joins.SetEntries(rc.Joins.Entries()...)
		sq = sq.Where(rc.Query, rc.Args...)
	}
	return RawCriteria{
		Joins: joins,
		Query: sq,
	}, nil
}

type NotCriteria struct {
	Criteria []Criteria
}

func (c NotCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	joins := maps.NewInsertMap[string, struct{}]()
	sq := ctx.Query().Torrent.UnderlyingDB()
	for _, cr := range c.Criteria {
		rc, rawCriteriaErr := cr.Raw(ctx)
		if rawCriteriaErr != nil {
			return RawCriteria{}, rawCriteriaErr
		}
		joins.SetEntries(rc.Joins.Entries()...)
		sq = sq.Not(rc.Query, rc.Args...)
	}
	return RawCriteria{
		Joins: joins,
		Query: sq,
	}, nil
}

type DbCriteria struct {
	Sql  string
	Args []interface{}
}

func (c DbCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	return RawCriteria{
		Query: c.Sql,
		Args:  c.Args,
	}, nil
}

type GenCriteria func(ctx DbContext) (Criteria, error)

func (c GenCriteria) Raw(ctx DbContext) (RawCriteria, error) {
	cc, err := c(ctx)
	if err != nil {
		return RawCriteria{}, err
	}
	return cc.Raw(ctx)
}

func queryStringCriteriaFromTokens(str string, tokens []string) Criteria {
	if len(tokens) == 0 {
		return OrCriteria{}
	}
	return GenCriteria(func(ctx DbContext) (Criteria, error) {
		return DbCriteria{
			Sql: strings.Join([]string{
				ctx.TableName() + ".tsv @@ plainto_tsquery('simple', ?)",
				ctx.TableName() + ".tsv @@ websearch_to_tsquery('simple', ?)",
				ctx.TableName() + ".search_string LIKE ?",
			}, " OR "),
			Args: []interface{}{
				strings.Join(tokens, " "),
				strings.Join(tokens, " "),
				"%" + strings.TrimSpace(str) + "%",
			},
		}, nil
	})
}

func QueryStringCriteria(str string) Criteria {
	queryStringTokens := regex.SearchStringToNormalizedTokens(str)
	return queryStringCriteriaFromTokens(str, queryStringTokens)
}
