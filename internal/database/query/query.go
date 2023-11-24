package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"sync"
)

type ResultItem struct {
	QueryStringRank float64
}

type GenericResult[T interface{}] struct {
	TotalCount   uint
	Items        []T
	Aggregations Aggregations
}

type SubQueryFactory = func(context.Context, *dao.Query) SubQuery

func GenericQuery[T interface{}](
	ctx context.Context,
	daoQ *dao.Query,
	option Option,
	tableName string,
	factory SubQueryFactory,
) (r GenericResult[T], _ error) {
	dbCtx := dbContext{
		q:         daoQ,
		tableName: tableName,
		factory:   factory,
	}
	builder, optionErr := option(newQueryContext(dbCtx))
	ctx = builder.createContext(ctx)
	if optionErr != nil {
		return r, optionErr
	}
	newSubQuery := func() (SubQuery, error) {
		sq := factory(ctx, daoQ)
		if selectErr := builder.applySelect(sq); selectErr != nil {
			return sq, selectErr
		}
		if preErr := builder.applyPre(sq); preErr != nil {
			return sq, preErr
		}
		return sq, nil
	}
	var errs []error
	mtx := sync.Mutex{}
	addErr := func(err error) {
		mtx.Lock()
		defer mtx.Unlock()
		errs = append(errs, err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go (func() {
		defer wg.Done()
		if builder.withTotalCount() {
			sq, sqErr := newSubQuery()
			if sqErr != nil {
				addErr(sqErr)
				return
			}
			if tc, countErr := sq.Count(); countErr != nil {
				addErr(countErr)
			} else {
				r.TotalCount = uint(tc)
			}
		}
	})()
	go (func() {
		defer wg.Done()
		sq, sqErr := newSubQuery()
		if sqErr != nil {
			addErr(sqErr)
			return
		}
		if postErr := builder.applyPost(sq); postErr != nil {
			addErr(postErr)
			return
		}
		if !builder.hasZeroLimit() {
			var items []T
			if txErr := sq.UnderlyingDB().Find(&items).Error; txErr != nil {
				addErr(txErr)
				return
			}
			if cbErr := builder.applyCallbacks(ctx, items); cbErr != nil {
				addErr(cbErr)
				return
			}
			r.Items = items
		}
	})()
	go (func() {
		defer wg.Done()
		if aggs, aggErr := builder.calculateAggregations(ctx); aggErr != nil {
			addErr(aggErr)
		} else {
			r.Aggregations = aggs
		}
	})()
	wg.Wait()
	return r, errors.Join(errs...)
}

type BaseSubQuery interface {
	Count() (int64, error)
	TableName() string
	UnderlyingDB() *gorm.DB
}

type SubQuery interface {
	BaseSubQuery
	Scopes(...GormScope) SubQuery
}

type GenericSubQuery[T BaseSubQuery] struct {
	SubQuery BaseSubQuery
}

func (sq GenericSubQuery[T]) Count() (int64, error) {
	return sq.SubQuery.Count()
}

func (sq GenericSubQuery[T]) TableName() string {
	return sq.SubQuery.TableName()
}

func (sq GenericSubQuery[T]) UnderlyingDB() *gorm.DB {
	return sq.SubQuery.UnderlyingDB()
}

type scoper[T BaseSubQuery] interface {
	Scopes(funcs ...GormScope) T
}

func (sq GenericSubQuery[T]) Scopes(fns ...func(gen.Dao) gen.Dao) SubQuery {
	sq.SubQuery = sq.SubQuery.(scoper[T]).Scopes(fns...)
	return sq
}

type Scope = func(SubQuery) error

type GormScope = func(gen.Dao) gen.Dao

type DbContext interface {
	Query() *dao.Query
	TableName() string
	NewSubQuery(context.Context) SubQuery
}

type dbContext struct {
	q         *dao.Query
	tableName string
	factory   SubQueryFactory
}

func (db dbContext) Query() *dao.Query {
	return db.q
}

func (db dbContext) TableName() string {
	return db.tableName
}

func (db dbContext) NewSubQuery(ctx context.Context) SubQuery {
	return db.factory(ctx, db.q)
}

type CallbackContext interface {
	DbContext
	Lock()
	Unlock()
}

type callbackContext struct {
	dbContext
	*sync.Mutex
}

type Callback func(ctx context.Context, cbCtx CallbackContext, results any) error

type OptionBuilder interface {
	DbContext
	Table(string) OptionBuilder
	Join(...TableJoin) OptionBuilder
	RequireJoin(...string) OptionBuilder
	Scope(...Scope) OptionBuilder
	GormScope(...GormScope) OptionBuilder
	Select(...clause.Expr) OptionBuilder
	OrderBy(...clause.OrderByColumn) OptionBuilder
	Limit(uint) OptionBuilder
	Offset(uint) OptionBuilder
	Group(...clause.Column) OptionBuilder
	Facet(...Facet) OptionBuilder
	Preload(...field.RelationField) OptionBuilder
	Callback(...Callback) OptionBuilder
	Context(func(ctx context.Context) context.Context) OptionBuilder
	applySelect(SubQuery) error
	applyPre(SubQuery) error
	applyPost(SubQuery) error
	createFacetsFilterCriteria() (Criteria, error)
	calculateAggregations(context.Context) (Aggregations, error)
	WithTotalCount(bool) OptionBuilder
	withTotalCount() bool
	applyCallbacks(context.Context, any) error
	hasZeroLimit() bool
	withCurrentFacet(string) OptionBuilder
	createContext(context.Context) context.Context
}

type optionBuilder struct {
	dbContext
	joins         map[string]TableJoin
	requiredJoins maps.InsertMap[string, struct{}]
	scopes        []Scope
	gormScopes    []GormScope
	selections    []clause.Expr
	groupBy       []clause.Column
	orderBy       []clause.OrderByColumn
	limit         model.NullUint
	offset        uint
	facets        []Facet
	currentFacet  string
	preloads      []field.RelationField
	totalCount    bool
	callbacks     []Callback
	contextFn     func(context.Context) context.Context
}

type RawJoin struct {
	Query string
	Args  []interface{}
}

func newQueryContext(dbCtx dbContext) OptionBuilder {
	return optionBuilder{
		dbContext:     dbCtx,
		joins:         make(map[string]TableJoin),
		requiredJoins: maps.NewInsertMap[string, struct{}](),
	}
}

func (b optionBuilder) Table(name string) OptionBuilder {
	b.dbContext.tableName = name
	return b.Scope(func(sq SubQuery) error {
		sq.UnderlyingDB().Table(name)
		return nil
	})
}

func (b optionBuilder) Join(joins ...TableJoin) OptionBuilder {
	bJoins := make(map[string]TableJoin, len(b.joins))
	for _, j := range b.joins {
		bJoins[j.Table.TableName()] = j
	}
	bRequiredJoins := b.requiredJoins.Copy()
	for _, j := range joins {
		bJoins[j.Table.TableName()] = j
		if j.Required {
			bRequiredJoins.SetKey(j.Table.TableName())
		}
	}
	b.joins = bJoins
	b.requiredJoins = bRequiredJoins
	return b
}

func (b optionBuilder) RequireJoin(names ...string) OptionBuilder {
	bRequiredJoins := b.requiredJoins.Copy()
	for _, name := range names {
		bRequiredJoins.SetKey(name)
	}
	b.requiredJoins = bRequiredJoins
	return b
}

func (b optionBuilder) Scope(scopes ...Scope) OptionBuilder {
	b.scopes = append(b.scopes, scopes...)
	return b
}

func (b optionBuilder) GormScope(scopes ...GormScope) OptionBuilder {
	b.gormScopes = append(b.gormScopes, scopes...)
	return b
}

func (b optionBuilder) Select(selections ...clause.Expr) OptionBuilder {
	b.selections = append(b.selections, selections...)
	return b
}

func (b optionBuilder) Group(columns ...clause.Column) OptionBuilder {
	b.groupBy = append(b.groupBy, columns...)
	return b
}

func (b optionBuilder) OrderBy(columns ...clause.OrderByColumn) OptionBuilder {
	b.orderBy = append(b.orderBy, columns...)
	return b
}

func (b optionBuilder) Limit(limit uint) OptionBuilder {
	b.limit = model.NewNullUint(limit)
	return b
}

func (b optionBuilder) Offset(offset uint) OptionBuilder {
	b.offset = offset
	return b
}

func (b optionBuilder) Facet(facets ...Facet) OptionBuilder {
	b.facets = append(b.facets, facets...)
	return b
}

func (b optionBuilder) Preload(relations ...field.RelationField) OptionBuilder {
	b.preloads = append(b.preloads, relations...)
	return b
}

func (b optionBuilder) Callback(callbacks ...Callback) OptionBuilder {
	b.callbacks = append(b.callbacks, callbacks...)
	return b
}

func (b optionBuilder) Context(fn func(context.Context) context.Context) OptionBuilder {
	prevFn := b.contextFn
	b.contextFn = func(ctx context.Context) context.Context {
		if prevFn != nil {
			ctx = prevFn(ctx)
		}
		return fn(ctx)
	}
	return b
}

func (b optionBuilder) WithTotalCount(bl bool) OptionBuilder {
	b.totalCount = bl
	return b
}

func (b optionBuilder) withTotalCount() bool {
	return b.totalCount
}

func (b optionBuilder) hasZeroLimit() bool {
	return b.limit.Valid && b.limit.Uint == 0
}

func (b optionBuilder) withCurrentFacet(facet string) OptionBuilder {
	b.currentFacet = facet
	return b
}

func (b optionBuilder) createContext(ctx context.Context) context.Context {
	if b.contextFn != nil {
		return b.contextFn(ctx)
	}
	return ctx
}

func (b optionBuilder) applySelect(sq SubQuery) error {
	var selectQueryParts []string
	selectQueryArgs := make([]interface{}, 0)
	if len(b.selections) == 0 {
		selectQueryParts = append(selectQueryParts, "*")
	} else {
		for _, s := range b.selections {
			selectQueryParts = append(selectQueryParts, s.SQL)
			selectQueryArgs = append(selectQueryArgs, s.Vars...)
		}
	}
	sq.UnderlyingDB().Select(strings.Join(selectQueryParts, ", "), selectQueryArgs...)
	return nil
}

func (b optionBuilder) applyPre(sq SubQuery) error {
	sq.Scopes(b.gormScopes...)
	for _, s := range b.scopes {
		if err := s(sq); err != nil {
			return err
		}
	}
	requiredJoins := b.requiredJoins.Copy()
	aggC, aggCErr := b.createFacetsFilterCriteria()
	if aggCErr != nil {
		return aggCErr
	}
	rawAggC, rawAggCErr := aggC.Raw(b)
	if rawAggCErr != nil {
		return rawAggCErr
	}
	requiredJoins.SetEntries(rawAggC.Joins.Entries()...)
	joins, joinsErr := extractRequiredJoins(b.dbContext.tableName, b.joins, requiredJoins)
	if joinsErr != nil {
		return joinsErr
	}
	applyJoins(sq, joins...)
	sq.UnderlyingDB().Where(rawAggC.Query, rawAggC.Args...)
	if len(b.groupBy) > 0 {
		sq.UnderlyingDB().Clauses(clause.GroupBy{
			Columns: b.groupBy,
		})
	}
	return nil
}

func extractRequiredJoins(tableName string, joins map[string]TableJoin, requiredJoins maps.InsertMap[string, struct{}]) ([]TableJoin, error) {
	resolvedJoins := maps.NewInsertMap[string, TableJoin]()
	var addJoin func(name string) error
	addJoin = func(name string) error {
		if name == tableName {
			return nil
		}
		j, ok := joins[name]
		if !ok {
			return fmt.Errorf("required join not found: %s", name)
		}
		for _, depName := range j.Dependencies.Keys() {
			if err := addJoin(depName); err != nil {
				return err
			}
		}
		resolvedJoins.Set(j.Table.TableName(), j)
		return nil
	}
	for _, joinName := range requiredJoins.Keys() {
		if err := addJoin(joinName); err != nil {
			return nil, err
		}
	}
	return resolvedJoins.Values(), nil
}

func applyJoins(sq SubQuery, joins ...TableJoin) {
	for _, j := range joins {
		join := j
		sq.Scopes(func(dao gen.Dao) gen.Dao {
			switch join.Type {
			case TableJoinTypeInner:
				return dao.Join(join.Table, join.On...)
			case TableJoinTypeLeft:
				return dao.LeftJoin(join.Table, join.On...)
			case TableJoinTypeRight:
				return dao.RightJoin(join.Table, join.On...)
			}
			panic("invalid join type")
		})
	}
}

func (b optionBuilder) applyPost(sq SubQuery) error {
	if len(b.orderBy) > 0 {
		sq.UnderlyingDB().Statement.AddClause(clause.OrderBy{
			Columns: b.orderBy,
		})
	}
	if b.limit.Valid {
		sq.UnderlyingDB().Limit(int(b.limit.Uint))
	}
	sq.UnderlyingDB().Offset(int(b.offset))
	for _, p := range b.preloads {
		sq.UnderlyingDB().Preload(p.Name(), p)
	}
	return nil
}

func (b optionBuilder) applyCallbacks(ctx context.Context, results any) error {
	cbCtx := callbackContext{
		dbContext: b.dbContext,
		Mutex:     &sync.Mutex{},
	}
	var errs []error
	wg := sync.WaitGroup{}
	wg.Add(len(b.callbacks))
	for _, cb := range b.callbacks {
		go (func(cb Callback) {
			defer wg.Done()
			if err := cb(ctx, cbCtx, results); err != nil {
				cbCtx.Lock()
				defer cbCtx.Unlock()
				errs = append(errs, err)
			}
		})(cb)
	}
	wg.Wait()
	return errors.Join(errs...)
}
