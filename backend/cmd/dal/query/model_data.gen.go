// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"backend/cmd/dal/model"
)

func newModelDatum(db *gorm.DB, opts ...gen.DOOption) modelDatum {
	_modelDatum := modelDatum{}

	_modelDatum.modelDatumDo.UseDB(db, opts...)
	_modelDatum.modelDatumDo.UseModel(&model.ModelDatum{})

	tableName := _modelDatum.modelDatumDo.TableName()
	_modelDatum.ALL = field.NewAsterisk(tableName)
	_modelDatum.ModelDataID = field.NewInt64(tableName, "model_data_id")
	_modelDatum.Data = field.NewString(tableName, "data")
	_modelDatum.ModelID = field.NewInt64(tableName, "model_id")
	_modelDatum.UserID = field.NewInt64(tableName, "user_id")

	_modelDatum.fillFieldMap()

	return _modelDatum
}

type modelDatum struct {
	modelDatumDo modelDatumDo

	ALL         field.Asterisk
	ModelDataID field.Int64  // 模型数据ID
	Data        field.String // 模型数据
	ModelID     field.Int64  // 模型ID
	UserID      field.Int64  // 用户ID

	fieldMap map[string]field.Expr
}

func (m modelDatum) Table(newTableName string) *modelDatum {
	m.modelDatumDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m modelDatum) As(alias string) *modelDatum {
	m.modelDatumDo.DO = *(m.modelDatumDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *modelDatum) updateTableName(table string) *modelDatum {
	m.ALL = field.NewAsterisk(table)
	m.ModelDataID = field.NewInt64(table, "model_data_id")
	m.Data = field.NewString(table, "data")
	m.ModelID = field.NewInt64(table, "model_id")
	m.UserID = field.NewInt64(table, "user_id")

	m.fillFieldMap()

	return m
}

func (m *modelDatum) WithContext(ctx context.Context) IModelDatumDo {
	return m.modelDatumDo.WithContext(ctx)
}

func (m modelDatum) TableName() string { return m.modelDatumDo.TableName() }

func (m modelDatum) Alias() string { return m.modelDatumDo.Alias() }

func (m *modelDatum) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *modelDatum) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 4)
	m.fieldMap["model_data_id"] = m.ModelDataID
	m.fieldMap["data"] = m.Data
	m.fieldMap["model_id"] = m.ModelID
	m.fieldMap["user_id"] = m.UserID
}

func (m modelDatum) clone(db *gorm.DB) modelDatum {
	m.modelDatumDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m modelDatum) replaceDB(db *gorm.DB) modelDatum {
	m.modelDatumDo.ReplaceDB(db)
	return m
}

type modelDatumDo struct{ gen.DO }

type IModelDatumDo interface {
	gen.SubQuery
	Debug() IModelDatumDo
	WithContext(ctx context.Context) IModelDatumDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IModelDatumDo
	WriteDB() IModelDatumDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IModelDatumDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IModelDatumDo
	Not(conds ...gen.Condition) IModelDatumDo
	Or(conds ...gen.Condition) IModelDatumDo
	Select(conds ...field.Expr) IModelDatumDo
	Where(conds ...gen.Condition) IModelDatumDo
	Order(conds ...field.Expr) IModelDatumDo
	Distinct(cols ...field.Expr) IModelDatumDo
	Omit(cols ...field.Expr) IModelDatumDo
	Join(table schema.Tabler, on ...field.Expr) IModelDatumDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IModelDatumDo
	RightJoin(table schema.Tabler, on ...field.Expr) IModelDatumDo
	Group(cols ...field.Expr) IModelDatumDo
	Having(conds ...gen.Condition) IModelDatumDo
	Limit(limit int) IModelDatumDo
	Offset(offset int) IModelDatumDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IModelDatumDo
	Unscoped() IModelDatumDo
	Create(values ...*model.ModelDatum) error
	CreateInBatches(values []*model.ModelDatum, batchSize int) error
	Save(values ...*model.ModelDatum) error
	First() (*model.ModelDatum, error)
	Take() (*model.ModelDatum, error)
	Last() (*model.ModelDatum, error)
	Find() ([]*model.ModelDatum, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ModelDatum, err error)
	FindInBatches(result *[]*model.ModelDatum, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ModelDatum) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IModelDatumDo
	Assign(attrs ...field.AssignExpr) IModelDatumDo
	Joins(fields ...field.RelationField) IModelDatumDo
	Preload(fields ...field.RelationField) IModelDatumDo
	FirstOrInit() (*model.ModelDatum, error)
	FirstOrCreate() (*model.ModelDatum, error)
	FindByPage(offset int, limit int) (result []*model.ModelDatum, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IModelDatumDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m modelDatumDo) Debug() IModelDatumDo {
	return m.withDO(m.DO.Debug())
}

func (m modelDatumDo) WithContext(ctx context.Context) IModelDatumDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m modelDatumDo) ReadDB() IModelDatumDo {
	return m.Clauses(dbresolver.Read)
}

func (m modelDatumDo) WriteDB() IModelDatumDo {
	return m.Clauses(dbresolver.Write)
}

func (m modelDatumDo) Session(config *gorm.Session) IModelDatumDo {
	return m.withDO(m.DO.Session(config))
}

func (m modelDatumDo) Clauses(conds ...clause.Expression) IModelDatumDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m modelDatumDo) Returning(value interface{}, columns ...string) IModelDatumDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m modelDatumDo) Not(conds ...gen.Condition) IModelDatumDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m modelDatumDo) Or(conds ...gen.Condition) IModelDatumDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m modelDatumDo) Select(conds ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m modelDatumDo) Where(conds ...gen.Condition) IModelDatumDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m modelDatumDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IModelDatumDo {
	return m.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (m modelDatumDo) Order(conds ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m modelDatumDo) Distinct(cols ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m modelDatumDo) Omit(cols ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m modelDatumDo) Join(table schema.Tabler, on ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m modelDatumDo) LeftJoin(table schema.Tabler, on ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m modelDatumDo) RightJoin(table schema.Tabler, on ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m modelDatumDo) Group(cols ...field.Expr) IModelDatumDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m modelDatumDo) Having(conds ...gen.Condition) IModelDatumDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m modelDatumDo) Limit(limit int) IModelDatumDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m modelDatumDo) Offset(offset int) IModelDatumDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m modelDatumDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IModelDatumDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m modelDatumDo) Unscoped() IModelDatumDo {
	return m.withDO(m.DO.Unscoped())
}

func (m modelDatumDo) Create(values ...*model.ModelDatum) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m modelDatumDo) CreateInBatches(values []*model.ModelDatum, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m modelDatumDo) Save(values ...*model.ModelDatum) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m modelDatumDo) First() (*model.ModelDatum, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModelDatum), nil
	}
}

func (m modelDatumDo) Take() (*model.ModelDatum, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModelDatum), nil
	}
}

func (m modelDatumDo) Last() (*model.ModelDatum, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModelDatum), nil
	}
}

func (m modelDatumDo) Find() ([]*model.ModelDatum, error) {
	result, err := m.DO.Find()
	return result.([]*model.ModelDatum), err
}

func (m modelDatumDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ModelDatum, err error) {
	buf := make([]*model.ModelDatum, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m modelDatumDo) FindInBatches(result *[]*model.ModelDatum, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m modelDatumDo) Attrs(attrs ...field.AssignExpr) IModelDatumDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m modelDatumDo) Assign(attrs ...field.AssignExpr) IModelDatumDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m modelDatumDo) Joins(fields ...field.RelationField) IModelDatumDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m modelDatumDo) Preload(fields ...field.RelationField) IModelDatumDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m modelDatumDo) FirstOrInit() (*model.ModelDatum, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModelDatum), nil
	}
}

func (m modelDatumDo) FirstOrCreate() (*model.ModelDatum, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModelDatum), nil
	}
}

func (m modelDatumDo) FindByPage(offset int, limit int) (result []*model.ModelDatum, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m modelDatumDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m modelDatumDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m modelDatumDo) Delete(models ...*model.ModelDatum) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *modelDatumDo) withDO(do gen.Dao) *modelDatumDo {
	m.DO = *do.(*gen.DO)
	return m
}
