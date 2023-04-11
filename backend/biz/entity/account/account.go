package account

import (
	"backend/biz/microtype"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
)

type Permission int

const (
	Normal  Permission = 1
	Manager Permission = 2
	Creator Permission = 3
)

type Operate int

const (
	Create    Operate = 1 // 创建账号
	IdQuery   Operate = 2 // 指定ID查询信息
	NameQuery Operate = 3 // 指定app_id, account_name查询信息
)

type Account struct {
	role Operate
	// Create
	appId       int64
	accountName string
	accountPwd  string
	permission  Permission
	// IdQuery
	accountId int64
	// NameQuery
	// appId       int64
	// accountName string
	
	addMo   *model.Account
	queryMo *model.Account
}

func NewAccount(accountId int64, appId int64, accountName, accountPwd string, permission Permission, role Operate) *Account {
	return &Account{
		accountId:   accountId,
		appId:       appId,
		permission:  permission,
		accountName: accountName,
		accountPwd:  accountPwd,
		role:        role,
	}
}

func (a *Account) Load(ctx context.Context) error {
	switch a.role {
	case Create:
		ac, err := query.Account.WithContext(ctx).Where(query.Account.AppID.Eq(a.appId), query.Account.AccountName.Eq(a.accountName)).First()
		if err != nil {
			return microtype.AccountQueryFailed
		}

		if ac != nil {
			return microtype.AccountExist
		}

		mo := &model.Account{
			AccountName:       a.accountName,
			AccountPwd:        a.accountPwd,
			AccountPermission: int64(a.permission),
			AppID:             a.appId,
		}

		err = query.Account.WithContext(ctx).Create(mo)
		if err != nil {
			return microtype.AccountAddFailed
		}

		a.addMo = mo
	case IdQuery:
		ac, err := query.Account.WithContext(ctx).Where(query.Account.AccountID.Eq(a.accountId)).First()
		if err != nil {
			return microtype.AccountQueryFailed
		}

		if ac == nil {
			return microtype.AccountNotExist
		}

		a.queryMo = ac
	case NameQuery:
		ac, err := query.Account.WithContext(ctx).Where(query.Account.AppID.Eq(a.appId), query.Account.AccountName.Eq(a.accountName)).First()
		if err != nil {
			return microtype.AccountQueryFailed
		}

		if ac == nil {
			return microtype.AccountNotExist
		}

		a.queryMo = ac
	default:
	}

	return nil
}

func (a *Account) GetQueryAccount() *model.Account {
	return a.queryMo
}

func (a *Account) GetAddedAccount() *model.Account {
	return a.addMo
}
