package crowd

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type Crowds struct {
	accountId int64

	//
	appId  int64
	total  int64
	crowds []*backend.Crowd
}

func NewCrowds(accountId int64) *Crowds {
	return &Crowds{
		accountId: accountId,
	}
}

func (c *Crowds) Load(ctx context.Context) error {
	ac := account.NewAccount(c.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	c.appId = ac.GetQueryAccount().AppID

	crowds, err := query.Crowd.WithContext(ctx).Where(query.Crowd.AppID.Eq(c.appId)).Find()
	if err != nil {
		return err
	}

	for _, crowd := range crowds {
		c.crowds = append(c.crowds, &backend.Crowd{
			CrowdName: crowd.CrowdName,
			CrowdDesc: crowd.CrowdDesc,
			CrowdID:   crowd.CrowdID,
		})
	}

	return nil
}

func (c *Crowds) GetResp() *backend.CrowdsResp {
	return &backend.CrowdsResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Total:      c.total,
		Crowds:     c.crowds,
	}
}
