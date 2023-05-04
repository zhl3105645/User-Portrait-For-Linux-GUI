package account

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
)

type AddAccount struct {
	accountId int64
	req       backend.AddAccountReq

	//
	appId             int64
	accountName       string
	accountPermission int64
	accountPwd        string
}

func NewAddAccount(accountId int64, req backend.AddAccountReq) *AddAccount {
	return &AddAccount{
		accountId: accountId,
		req:       req,
	}
}

func (a *AddAccount) Load(ctx context.Context) error {
	if !a.check() {
		return microtype.ParamCheckFailed
	}
	ac := account.NewAccount(a.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}
	a.appId = ac.GetQueryAccount().AppID

	mo := &model.Account{
		AccountID:         0,
		AccountName:       a.accountName,
		AccountPwd:        a.accountPwd,
		AccountPermission: a.accountPermission,
		AppID:             a.appId,
	}

	err := query.Account.WithContext(ctx).Create(mo)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddAccount) check() bool {
	if a.req.AccountName == "" || a.req.AccountPermission <= 0 || a.req.AccountPwd == "" {
		return false
	}

	a.accountName = a.req.AccountName
	a.accountPermission = a.req.AccountPermission
	a.accountPwd = a.req.AccountPwd

	return true
}

func (a *AddAccount) GetResp() *backend.AddAccountResp {
	return &backend.AddAccountResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
