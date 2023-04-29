package crowd

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
)

type AddCrowd struct {
	accountId int64
	req       backend.AddCrowdReq

	//
	appId int64
}

func NewAddCrowd(accountId int64, req backend.AddCrowdReq) *AddCrowd {
	return &AddCrowd{
		accountId: accountId,
		req:       req,
	}
}

func (a *AddCrowd) Load(ctx context.Context) error {
	if !a.check() {
		return microtype.ParamCheckFailed
	}

	ac := account.NewAccount(a.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}
	a.appId = ac.GetQueryAccount().AppID

	rules, err := json.Marshal(a.req.DivideRules)
	if err != nil {
		return microtype.JsonMarshalFailed
	}

	mo := &model.Crowd{
		CrowdID:         0,
		CrowdDesc:       a.req.CrowdDesc,
		AppID:           a.appId,
		CrowdName:       a.req.CrowdName,
		CrowdDivideRule: string(rules),
	}

	err = query.Crowd.WithContext(ctx).Create(mo)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddCrowd) GetResp() *backend.AddCrowdResp {
	return &backend.AddCrowdResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

func (a *AddCrowd) check() bool {
	if a.req.CrowdName == "" || a.req.CrowdDesc == "" || len(a.req.DivideRules) == 0 {
		return false
	}

	return true
}
