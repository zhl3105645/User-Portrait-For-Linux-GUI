package account

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type PageAccount struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string
	//
	appId    int64
	total    int64
	accounts []*backend.Account
}

func NewPageAccount(accountId int64, pageNum int64, pageSize int64, search string) *PageAccount {
	return &PageAccount{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

func (p *PageAccount) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	offset := (p.pageNum - 1) * p.pageSize
	res, count, err := query.Account.WithContext(ctx).
		Where(query.Account.AppID.Eq(p.appId),
			query.Account.AccountID.NotIn(p.accountId),
			query.Account.AccountName.Like("%"+p.search+"%")).
		FindByPage(int(offset), int(p.pageSize))
	if err != nil {
		return err
	}
	p.total = count
	p.accounts = make([]*backend.Account, 0, len(res))
	for _, r := range res {
		if r == nil {
			continue
		}
		p.accounts = append(p.accounts, &backend.Account{
			AccountName:       r.AccountName,
			AccountPermission: r.AccountPermission,
			AccountID:         r.AccountID,
		})
	}

	return nil
}

func (p *PageAccount) GetResp() *backend.AccountInPageResp {
	return &backend.AccountInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Accounts:   p.accounts,
		Total:      p.total,
	}
}
