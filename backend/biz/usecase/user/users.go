package user

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Users struct {
	accountId int64

	//
	appId int64
	res   []*backend.User
}

func NewUsers(accountId int64) *Users {
	return &Users{
		accountId: accountId,
	}
}

func (u *Users) Load(ctx context.Context) error {
	ac := account.NewAccount(u.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}
	u.appId = ac.GetQueryAccount().AppID

	users, err := query.User.WithContext(ctx).Where(query.User.AppID.Eq(u.appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return microtype.UserQueryFailed
	}

	u.res = make([]*backend.User, 0, len(users))
	for _, user := range users {
		if user == nil {
			continue
		}
		u.res = append(u.res, &backend.User{
			UserID:   user.UserID,
			UserName: user.UserName,
		})
	}

	return nil
}

func (u *Users) GetResp() *backend.UsersResp {
	return &backend.UsersResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Users:      u.res,
	}
}
