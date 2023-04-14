package user

import (
	"backend/biz/entity/account"
	"backend/biz/entity/user"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
)

type User struct {
	accountId int64
	userName  string
}

func NewUser(accountId int64, userName string) *User {
	return &User{
		accountId: accountId,
		userName:  userName,
	}
}

func (u *User) Load(ctx context.Context) error {
	if u.userName == "" {
		return microtype.UserNameEmpty
	}
	ac := account.NewAccount(u.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	ue := user.NewUser(user.Create, &user.CreateParam{
		AppId:    ac.GetQueryAccount().AppID,
		UserName: u.userName,
	}, nil, 0)
	if err := ue.Load(ctx); err != nil {
		return err
	}

	return nil
}

func (u *User) GetResp() *backend.AddUserResp {
	return &backend.AddUserResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
