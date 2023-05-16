package rule

import (
	"backend/biz/entity/account"
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
)

type AddRule struct {
	accountId int64
	req       backend.AddRuleReq

	//
	appId int64
}

func NewAddRule(accountId int64, req backend.AddRuleReq) *AddRule {
	return &AddRule{
		accountId: accountId,
		req:       req,
	}
}

func (r *AddRule) Load(ctx context.Context) error {
	if r.req.RuleType != int64(rule.EventRule) && r.req.RuleType != int64(rule.BehaviorRule) {
		return microtype.RuleParamFailed
	}

	if r.req.RuleDesc == "" {
		return microtype.ParamCheckFailed
	}

	ac := account.NewAccount(r.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	r.appId = ac.GetQueryAccount().AppID

	do := query.Rule
	mo := do.WithContext(ctx)

	queryMo, err := mo.Where(do.AppID.Eq(r.appId),
		do.RuleType.Eq(r.req.RuleType),
		do.RuleDesc.Eq(r.req.RuleDesc)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("rule query failed. err=", err.Error())
		return microtype.RuleQueryFailed
	}

	if queryMo != nil {
		return microtype.RuleExist
	}

	createMo := &model.Rule{
		RuleID:   0,
		RuleType: r.req.RuleType,
		RuleDesc: r.req.RuleDesc,
		AppID:    r.appId,
	}

	err = query.Rule.WithContext(ctx).Create(createMo)
	if err != nil {
		logger.Error("rule create failed. err=", err.Error())
		return microtype.RuleCreateFailed
	}

	return nil
}

func (r *AddRule) GetResp() *backend.AddRuleResp {
	return &backend.AddRuleResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
