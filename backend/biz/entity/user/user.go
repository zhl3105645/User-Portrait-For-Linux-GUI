package user

import (
	"backend/biz/microtype"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Operate int

const (
	Create    Operate = 1 // 创建用户
	QueryPage Operate = 2 // 分页查询
	IdQuery   Operate = 3 // 指定ID查询
)

type User struct {
	role Operate
	// Create
	createParam *CreateParam
	// QueryPage
	queryParam *QueryParam
	// IdQuery
	userId int64

	// Create
	addMo *model.User
	// QueryPage
	total  int64
	pageMo []*model.User
	// IdQuery
	queryMo *model.User
}

type CreateParam struct {
	AppId    int64
	UserName string
}

type QueryParam struct {
	AppId    int64
	PageNum  int64
	PageSize int64
	Search   string
}

func NewUser(role Operate, createParam *CreateParam, queryParam *QueryParam, userId int64) *User {
	return &User{
		role:        role,
		createParam: createParam,
		queryParam:  queryParam,
		userId:      userId,
	}
}

func (u *User) Load(ctx context.Context) error {
	switch u.role {
	case Create:
		uQuery, err := query.User.WithContext(ctx).
			Where(query.User.UserName.Eq(u.createParam.UserName), query.User.AppID.Eq(u.createParam.AppId)).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return microtype.UserQueryFailed
		}

		if uQuery != nil {
			return microtype.UserExist
		}

		mo := &model.User{
			UserName: u.createParam.UserName,
			AppID:    u.createParam.AppId,
		}

		err = query.User.WithContext(ctx).Create(mo)
		if err != nil {
			return microtype.UserAddFailed
		}

		u.addMo = mo
	case QueryPage:
		// 分页
		offset := (u.queryParam.PageNum - 1) * u.queryParam.PageSize
		res, count, err := query.User.WithContext(ctx).
			Where(query.User.AppID.Eq(u.queryParam.AppId), query.User.UserName.Like("%"+u.queryParam.Search+"%")).
			FindByPage(int(offset), int(u.queryParam.PageSize))
		if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
			return microtype.UserQueryFailed
		}

		u.total = count
		u.pageMo = res
	case IdQuery:
		queryMo, err := query.User.WithContext(ctx).
			Where(query.User.UserID.Eq(u.userId)).
			First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return microtype.UserQueryFailed
		}

		u.queryMo = queryMo
	default:
	}

	return nil
}

func (u *User) GetAddedUser() *model.User {
	return u.addMo
}

func (u *User) GetTotal() int64 {
	return u.total
}

func (u *User) GetPageQueryUser() []*model.User {
	return u.pageMo
}

func (u *User) GetIdQueryUser() *model.User {
	return u.queryMo
}
