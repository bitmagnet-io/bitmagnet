package query

import (
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func Where(conditions ...Criteria) Option {
	if len(conditions) == 0 {
		return Options()
	}

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
	//revive:disable-next-line:nested-structs
	Joins maps.InsertMap[string, struct{}]
}

func (c RawCriteria) Raw(DBContext) (RawCriteria, error) {
	return c, nil
}

type DaoCriteria struct {
	Conditions func(ctx DBContext) ([]field.Expr, error)
	//revive:disable-next-line:nested-structs
	Joins maps.InsertMap[string, struct{}]
}

func (c DaoCriteria) Raw(ctx DBContext) (RawCriteria, error) {
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

func (c OrCriteria) Raw(ctx DBContext) (RawCriteria, error) {
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

func (c AndCriteria) Raw(ctx DBContext) (RawCriteria, error) {
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

func (c NotCriteria) Raw(ctx DBContext) (RawCriteria, error) {
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

type DBCriteria struct {
	SQL  string
	Args []interface{}
}

func (c DBCriteria) Raw(DBContext) (RawCriteria, error) {
	return RawCriteria{
		Query: c.SQL,
		Args:  c.Args,
	}, nil
}

type GenCriteria func(ctx DBContext) (Criteria, error)

func (c GenCriteria) Raw(ctx DBContext) (RawCriteria, error) {
	cc, err := c(ctx)
	if err != nil {
		return RawCriteria{}, err
	}

	return cc.Raw(ctx)
}
