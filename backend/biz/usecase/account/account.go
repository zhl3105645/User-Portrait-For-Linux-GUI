package account

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"context"
)

type Account struct {
	accountId int64

	queryAccount *model.Account
}

func NewAccount(accountId int64) *Account {
	return &Account{
		accountId: accountId,
	}
}

func (a *Account) Load(ctx context.Context) error {
	ac := account.NewAccount(a.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	a.queryAccount = ac.GetQueryAccount()

	return nil
}

func (a *Account) GetResp() *backend.AccountResp {
	return &backend.AccountResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Account: &backend.Account{
			AccountName:       a.queryAccount.AccountName,
			AccountPermission: a.queryAccount.AccountPermission,
		},
	}
}
