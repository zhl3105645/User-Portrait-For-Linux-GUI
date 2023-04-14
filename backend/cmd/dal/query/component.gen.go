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

func newComponent(db *gorm.DB, opts ...gen.DOOption) component {
	_component := component{}

	_component.componentDo.UseDB(db, opts...)
	_component.componentDo.UseModel(&model.Component{})

	tableName := _component.componentDo.TableName()
	_component.ALL = field.NewAsterisk(tableName)
	_component.ComponentID = field.NewInt64(tableName, "component_id")
	_component.ComponentName = field.NewString(tableName, "component_name")
	_component.ComponentType = field.NewInt64(tableName, "component_type")
	_component.AppID = field.NewInt64(tableName, "app_id")
	_component.ComponentDesc = field.NewString(tableName, "component_desc")

	_component.fillFieldMap()

	return _component
}

type component struct {
	componentDo componentDo

	ALL           field.Asterisk
	ComponentID   field.Int64  // 组件ID
	ComponentName field.String // 组件名
	ComponentType field.Int64  // 组件类型
	AppID         field.Int64  // 应用ID
	ComponentDesc field.String // 组件描述

	fieldMap map[string]field.Expr
}

func (c component) Table(newTableName string) *component {
	c.componentDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c component) As(alias string) *component {
	c.componentDo.DO = *(c.componentDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *component) updateTableName(table string) *component {
	c.ALL = field.NewAsterisk(table)
	c.ComponentID = field.NewInt64(table, "component_id")
	c.ComponentName = field.NewString(table, "component_name")
	c.ComponentType = field.NewInt64(table, "component_type")
	c.AppID = field.NewInt64(table, "app_id")
	c.ComponentDesc = field.NewString(table, "component_desc")

	c.fillFieldMap()

	return c
}

func (c *component) WithContext(ctx context.Context) IComponentDo {
	return c.componentDo.WithContext(ctx)
}

func (c component) TableName() string { return c.componentDo.TableName() }

func (c component) Alias() string { return c.componentDo.Alias() }

func (c *component) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *component) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 5)
	c.fieldMap["component_id"] = c.ComponentID
	c.fieldMap["component_name"] = c.ComponentName
	c.fieldMap["component_type"] = c.ComponentType
	c.fieldMap["app_id"] = c.AppID
	c.fieldMap["component_desc"] = c.ComponentDesc
}

func (c component) clone(db *gorm.DB) component {
	c.componentDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c component) replaceDB(db *gorm.DB) component {
	c.componentDo.ReplaceDB(db)
	return c
}

type componentDo struct{ gen.DO }

type IComponentDo interface {
	gen.SubQuery
	Debug() IComponentDo
	WithContext(ctx context.Context) IComponentDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IComponentDo
	WriteDB() IComponentDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IComponentDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IComponentDo
	Not(conds ...gen.Condition) IComponentDo
	Or(conds ...gen.Condition) IComponentDo
	Select(conds ...field.Expr) IComponentDo
	Where(conds ...gen.Condition) IComponentDo
	Order(conds ...field.Expr) IComponentDo
	Distinct(cols ...field.Expr) IComponentDo
	Omit(cols ...field.Expr) IComponentDo
	Join(table schema.Tabler, on ...field.Expr) IComponentDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IComponentDo
	RightJoin(table schema.Tabler, on ...field.Expr) IComponentDo
	Group(cols ...field.Expr) IComponentDo
	Having(conds ...gen.Condition) IComponentDo
	Limit(limit int) IComponentDo
	Offset(offset int) IComponentDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IComponentDo
	Unscoped() IComponentDo
	Create(values ...*model.Component) error
	CreateInBatches(values []*model.Component, batchSize int) error
	Save(values ...*model.Component) error
	First() (*model.Component, error)
	Take() (*model.Component, error)
	Last() (*model.Component, error)
	Find() ([]*model.Component, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Component, err error)
	FindInBatches(result *[]*model.Component, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Component) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IComponentDo
	Assign(attrs ...field.AssignExpr) IComponentDo
	Joins(fields ...field.RelationField) IComponentDo
	Preload(fields ...field.RelationField) IComponentDo
	FirstOrInit() (*model.Component, error)
	FirstOrCreate() (*model.Component, error)
	FindByPage(offset int, limit int) (result []*model.Component, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IComponentDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c componentDo) Debug() IComponentDo {
	return c.withDO(c.DO.Debug())
}

func (c componentDo) WithContext(ctx context.Context) IComponentDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c componentDo) ReadDB() IComponentDo {
	return c.Clauses(dbresolver.Read)
}

func (c componentDo) WriteDB() IComponentDo {
	return c.Clauses(dbresolver.Write)
}

func (c componentDo) Session(config *gorm.Session) IComponentDo {
	return c.withDO(c.DO.Session(config))
}

func (c componentDo) Clauses(conds ...clause.Expression) IComponentDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c componentDo) Returning(value interface{}, columns ...string) IComponentDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c componentDo) Not(conds ...gen.Condition) IComponentDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c componentDo) Or(conds ...gen.Condition) IComponentDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c componentDo) Select(conds ...field.Expr) IComponentDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c componentDo) Where(conds ...gen.Condition) IComponentDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c componentDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IComponentDo {
	return c.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (c componentDo) Order(conds ...field.Expr) IComponentDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c componentDo) Distinct(cols ...field.Expr) IComponentDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c componentDo) Omit(cols ...field.Expr) IComponentDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c componentDo) Join(table schema.Tabler, on ...field.Expr) IComponentDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c componentDo) LeftJoin(table schema.Tabler, on ...field.Expr) IComponentDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c componentDo) RightJoin(table schema.Tabler, on ...field.Expr) IComponentDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c componentDo) Group(cols ...field.Expr) IComponentDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c componentDo) Having(conds ...gen.Condition) IComponentDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c componentDo) Limit(limit int) IComponentDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c componentDo) Offset(offset int) IComponentDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c componentDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IComponentDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c componentDo) Unscoped() IComponentDo {
	return c.withDO(c.DO.Unscoped())
}

func (c componentDo) Create(values ...*model.Component) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c componentDo) CreateInBatches(values []*model.Component, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c componentDo) Save(values ...*model.Component) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c componentDo) First() (*model.Component, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Component), nil
	}
}

func (c componentDo) Take() (*model.Component, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Component), nil
	}
}

func (c componentDo) Last() (*model.Component, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Component), nil
	}
}

func (c componentDo) Find() ([]*model.Component, error) {
	result, err := c.DO.Find()
	return result.([]*model.Component), err
}

func (c componentDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Component, err error) {
	buf := make([]*model.Component, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c componentDo) FindInBatches(result *[]*model.Component, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c componentDo) Attrs(attrs ...field.AssignExpr) IComponentDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c componentDo) Assign(attrs ...field.AssignExpr) IComponentDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c componentDo) Joins(fields ...field.RelationField) IComponentDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c componentDo) Preload(fields ...field.RelationField) IComponentDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c componentDo) FirstOrInit() (*model.Component, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Component), nil
	}
}

func (c componentDo) FirstOrCreate() (*model.Component, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Component), nil
	}
}

func (c componentDo) FindByPage(offset int, limit int) (result []*model.Component, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c componentDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c componentDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c componentDo) Delete(models ...*model.Component) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *componentDo) withDO(do gen.Dao) *componentDo {
	c.DO = *do.(*gen.DO)
	return c
}
