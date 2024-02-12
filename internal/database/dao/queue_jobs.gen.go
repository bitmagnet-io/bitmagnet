// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func newQueueJob(db *gorm.DB, opts ...gen.DOOption) queueJob {
	_queueJob := queueJob{}

	_queueJob.queueJobDo.UseDB(db, opts...)
	_queueJob.queueJobDo.UseModel(&model.QueueJob{})

	tableName := _queueJob.queueJobDo.TableName()
	_queueJob.ALL = field.NewAsterisk(tableName)
	_queueJob.ID = field.NewString(tableName, "id")
	_queueJob.Fingerprint = field.NewString(tableName, "fingerprint")
	_queueJob.Queue = field.NewString(tableName, "queue")
	_queueJob.Status = field.NewString(tableName, "status")
	_queueJob.Payload = field.NewString(tableName, "payload")
	_queueJob.Retries = field.NewUint(tableName, "retries")
	_queueJob.MaxRetries = field.NewUint(tableName, "max_retries")
	_queueJob.RunAfter = field.NewTime(tableName, "run_after")
	_queueJob.RanAt = field.NewField(tableName, "ran_at")
	_queueJob.Error = field.NewField(tableName, "error")
	_queueJob.Deadline = field.NewField(tableName, "deadline")
	_queueJob.ArchivalDuration = field.NewField(tableName, "archival_duration")
	_queueJob.CreatedAt = field.NewTime(tableName, "created_at")

	_queueJob.fillFieldMap()

	return _queueJob
}

type queueJob struct {
	queueJobDo

	ALL              field.Asterisk
	ID               field.String
	Fingerprint      field.String
	Queue            field.String
	Status           field.String
	Payload          field.String
	Retries          field.Uint
	MaxRetries       field.Uint
	RunAfter         field.Time
	RanAt            field.Field
	Error            field.Field
	Deadline         field.Field
	ArchivalDuration field.Field
	CreatedAt        field.Time

	fieldMap map[string]field.Expr
}

func (q queueJob) Table(newTableName string) *queueJob {
	q.queueJobDo.UseTable(newTableName)
	return q.updateTableName(newTableName)
}

func (q queueJob) As(alias string) *queueJob {
	q.queueJobDo.DO = *(q.queueJobDo.As(alias).(*gen.DO))
	return q.updateTableName(alias)
}

func (q *queueJob) updateTableName(table string) *queueJob {
	q.ALL = field.NewAsterisk(table)
	q.ID = field.NewString(table, "id")
	q.Fingerprint = field.NewString(table, "fingerprint")
	q.Queue = field.NewString(table, "queue")
	q.Status = field.NewString(table, "status")
	q.Payload = field.NewString(table, "payload")
	q.Retries = field.NewUint(table, "retries")
	q.MaxRetries = field.NewUint(table, "max_retries")
	q.RunAfter = field.NewTime(table, "run_after")
	q.RanAt = field.NewField(table, "ran_at")
	q.Error = field.NewField(table, "error")
	q.Deadline = field.NewField(table, "deadline")
	q.ArchivalDuration = field.NewField(table, "archival_duration")
	q.CreatedAt = field.NewTime(table, "created_at")

	q.fillFieldMap()

	return q
}

func (q *queueJob) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := q.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (q *queueJob) fillFieldMap() {
	q.fieldMap = make(map[string]field.Expr, 13)
	q.fieldMap["id"] = q.ID
	q.fieldMap["fingerprint"] = q.Fingerprint
	q.fieldMap["queue"] = q.Queue
	q.fieldMap["status"] = q.Status
	q.fieldMap["payload"] = q.Payload
	q.fieldMap["retries"] = q.Retries
	q.fieldMap["max_retries"] = q.MaxRetries
	q.fieldMap["run_after"] = q.RunAfter
	q.fieldMap["ran_at"] = q.RanAt
	q.fieldMap["error"] = q.Error
	q.fieldMap["deadline"] = q.Deadline
	q.fieldMap["archival_duration"] = q.ArchivalDuration
	q.fieldMap["created_at"] = q.CreatedAt
}

func (q queueJob) clone(db *gorm.DB) queueJob {
	q.queueJobDo.ReplaceConnPool(db.Statement.ConnPool)
	return q
}

func (q queueJob) replaceDB(db *gorm.DB) queueJob {
	q.queueJobDo.ReplaceDB(db)
	return q
}

type queueJobDo struct{ gen.DO }

type IQueueJobDo interface {
	gen.SubQuery
	Debug() IQueueJobDo
	WithContext(ctx context.Context) IQueueJobDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IQueueJobDo
	WriteDB() IQueueJobDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IQueueJobDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IQueueJobDo
	Not(conds ...gen.Condition) IQueueJobDo
	Or(conds ...gen.Condition) IQueueJobDo
	Select(conds ...field.Expr) IQueueJobDo
	Where(conds ...gen.Condition) IQueueJobDo
	Order(conds ...field.Expr) IQueueJobDo
	Distinct(cols ...field.Expr) IQueueJobDo
	Omit(cols ...field.Expr) IQueueJobDo
	Join(table schema.Tabler, on ...field.Expr) IQueueJobDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IQueueJobDo
	RightJoin(table schema.Tabler, on ...field.Expr) IQueueJobDo
	Group(cols ...field.Expr) IQueueJobDo
	Having(conds ...gen.Condition) IQueueJobDo
	Limit(limit int) IQueueJobDo
	Offset(offset int) IQueueJobDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IQueueJobDo
	Unscoped() IQueueJobDo
	Create(values ...*model.QueueJob) error
	CreateInBatches(values []*model.QueueJob, batchSize int) error
	Save(values ...*model.QueueJob) error
	First() (*model.QueueJob, error)
	Take() (*model.QueueJob, error)
	Last() (*model.QueueJob, error)
	Find() ([]*model.QueueJob, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.QueueJob, err error)
	FindInBatches(result *[]*model.QueueJob, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.QueueJob) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IQueueJobDo
	Assign(attrs ...field.AssignExpr) IQueueJobDo
	Joins(fields ...field.RelationField) IQueueJobDo
	Preload(fields ...field.RelationField) IQueueJobDo
	FirstOrInit() (*model.QueueJob, error)
	FirstOrCreate() (*model.QueueJob, error)
	FindByPage(offset int, limit int) (result []*model.QueueJob, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IQueueJobDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (q queueJobDo) Debug() IQueueJobDo {
	return q.withDO(q.DO.Debug())
}

func (q queueJobDo) WithContext(ctx context.Context) IQueueJobDo {
	return q.withDO(q.DO.WithContext(ctx))
}

func (q queueJobDo) ReadDB() IQueueJobDo {
	return q.Clauses(dbresolver.Read)
}

func (q queueJobDo) WriteDB() IQueueJobDo {
	return q.Clauses(dbresolver.Write)
}

func (q queueJobDo) Session(config *gorm.Session) IQueueJobDo {
	return q.withDO(q.DO.Session(config))
}

func (q queueJobDo) Clauses(conds ...clause.Expression) IQueueJobDo {
	return q.withDO(q.DO.Clauses(conds...))
}

func (q queueJobDo) Returning(value interface{}, columns ...string) IQueueJobDo {
	return q.withDO(q.DO.Returning(value, columns...))
}

func (q queueJobDo) Not(conds ...gen.Condition) IQueueJobDo {
	return q.withDO(q.DO.Not(conds...))
}

func (q queueJobDo) Or(conds ...gen.Condition) IQueueJobDo {
	return q.withDO(q.DO.Or(conds...))
}

func (q queueJobDo) Select(conds ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.Select(conds...))
}

func (q queueJobDo) Where(conds ...gen.Condition) IQueueJobDo {
	return q.withDO(q.DO.Where(conds...))
}

func (q queueJobDo) Order(conds ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.Order(conds...))
}

func (q queueJobDo) Distinct(cols ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.Distinct(cols...))
}

func (q queueJobDo) Omit(cols ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.Omit(cols...))
}

func (q queueJobDo) Join(table schema.Tabler, on ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.Join(table, on...))
}

func (q queueJobDo) LeftJoin(table schema.Tabler, on ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.LeftJoin(table, on...))
}

func (q queueJobDo) RightJoin(table schema.Tabler, on ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.RightJoin(table, on...))
}

func (q queueJobDo) Group(cols ...field.Expr) IQueueJobDo {
	return q.withDO(q.DO.Group(cols...))
}

func (q queueJobDo) Having(conds ...gen.Condition) IQueueJobDo {
	return q.withDO(q.DO.Having(conds...))
}

func (q queueJobDo) Limit(limit int) IQueueJobDo {
	return q.withDO(q.DO.Limit(limit))
}

func (q queueJobDo) Offset(offset int) IQueueJobDo {
	return q.withDO(q.DO.Offset(offset))
}

func (q queueJobDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IQueueJobDo {
	return q.withDO(q.DO.Scopes(funcs...))
}

func (q queueJobDo) Unscoped() IQueueJobDo {
	return q.withDO(q.DO.Unscoped())
}

func (q queueJobDo) Create(values ...*model.QueueJob) error {
	if len(values) == 0 {
		return nil
	}
	return q.DO.Create(values)
}

func (q queueJobDo) CreateInBatches(values []*model.QueueJob, batchSize int) error {
	return q.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (q queueJobDo) Save(values ...*model.QueueJob) error {
	if len(values) == 0 {
		return nil
	}
	return q.DO.Save(values)
}

func (q queueJobDo) First() (*model.QueueJob, error) {
	if result, err := q.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.QueueJob), nil
	}
}

func (q queueJobDo) Take() (*model.QueueJob, error) {
	if result, err := q.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.QueueJob), nil
	}
}

func (q queueJobDo) Last() (*model.QueueJob, error) {
	if result, err := q.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.QueueJob), nil
	}
}

func (q queueJobDo) Find() ([]*model.QueueJob, error) {
	result, err := q.DO.Find()
	return result.([]*model.QueueJob), err
}

func (q queueJobDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.QueueJob, err error) {
	buf := make([]*model.QueueJob, 0, batchSize)
	err = q.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (q queueJobDo) FindInBatches(result *[]*model.QueueJob, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return q.DO.FindInBatches(result, batchSize, fc)
}

func (q queueJobDo) Attrs(attrs ...field.AssignExpr) IQueueJobDo {
	return q.withDO(q.DO.Attrs(attrs...))
}

func (q queueJobDo) Assign(attrs ...field.AssignExpr) IQueueJobDo {
	return q.withDO(q.DO.Assign(attrs...))
}

func (q queueJobDo) Joins(fields ...field.RelationField) IQueueJobDo {
	for _, _f := range fields {
		q = *q.withDO(q.DO.Joins(_f))
	}
	return &q
}

func (q queueJobDo) Preload(fields ...field.RelationField) IQueueJobDo {
	for _, _f := range fields {
		q = *q.withDO(q.DO.Preload(_f))
	}
	return &q
}

func (q queueJobDo) FirstOrInit() (*model.QueueJob, error) {
	if result, err := q.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.QueueJob), nil
	}
}

func (q queueJobDo) FirstOrCreate() (*model.QueueJob, error) {
	if result, err := q.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.QueueJob), nil
	}
}

func (q queueJobDo) FindByPage(offset int, limit int) (result []*model.QueueJob, count int64, err error) {
	result, err = q.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = q.Offset(-1).Limit(-1).Count()
	return
}

func (q queueJobDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = q.Count()
	if err != nil {
		return
	}

	err = q.Offset(offset).Limit(limit).Scan(result)
	return
}

func (q queueJobDo) Scan(result interface{}) (err error) {
	return q.DO.Scan(result)
}

func (q queueJobDo) Delete(models ...*model.QueueJob) (result gen.ResultInfo, err error) {
	return q.DO.Delete(models)
}

func (q *queueJobDo) withDO(do gen.Dao) *queueJobDo {
	q.DO = *do.(*gen.DO)
	return q
}
