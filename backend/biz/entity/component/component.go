package component

import (
	"backend/biz/microtype"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
)

type Operate int

const (
	QueryPage   Operate = 1 // 分页查询
	QueryAll    Operate = 2 // 全量查询
	InsertBatch Operate = 3 // 批量插入
)

type Component struct {
	role  Operate
	appId int64
	// QueryPage
	queryParam *QueryParam
	// QueryAll

	// InsertBatch
	insetParam *InsertParam

	// QueryPage
	total   int64
	queryMo []*model.Component
	// AllPage
	// queryMo []*model.Component
}

type QueryParam struct {
	PageNum  int64
	PageSize int64
	Search   string
}

type InsertParam struct {
	InsertMo []*model.Component
}

func NewComponent(role Operate, appId int64, queryParam *QueryParam, insetParam *InsertParam) *Component {
	return &Component{
		appId:      appId,
		role:       role,
		queryParam: queryParam,
		insetParam: insetParam,
	}
}

func (c *Component) Load(ctx context.Context) error {
	switch c.role {
	case QueryPage:
		offset := (c.queryParam.PageNum - 1) * c.queryParam.PageSize
		co := query.Component.WithContext(ctx)
		res, pageCount, err := co.Where(
			query.Component.AppID.Eq(c.appId),
			co.Where(query.Component.ComponentName.Like("%"+c.queryParam.Search+"%")).
				Or(query.Component.ComponentDesc.Like("%"+c.queryParam.Search+"%"))).
			FindByPage(int(offset), int(c.queryParam.PageSize))

		if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
			return microtype.ComponentQueryFailed
		}

		c.total = pageCount
		c.queryMo = res
	case QueryAll:
		res, err := query.Component.WithContext(ctx).
			Where(query.Component.AppID.Eq(c.appId)).
			Find()
		if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
			return microtype.ComponentQueryFailed
		}

		c.queryMo = res
	case InsertBatch:
		err := query.Component.WithContext(ctx).Create(c.insetParam.InsertMo...)
		if err != nil {
			logger.Error("insert component failed. err=", err.Error())
			return microtype.ComponentCreateFailed
		}
	default:
	}

	return nil
}

func (c *Component) GetTotal() int64 {
	return c.total
}

func (c *Component) GetQueryComponent() []*model.Component {
	return c.queryMo
}
