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
	"gorm.io/gorm"
)

type Rules struct {
	accountId int64
	ruleType  int64

	//
	appId int64
	res   []*model.Rule
}

func NewRules(accountId int64, ruleType int64) *Rules {
	return &Rules{
		accountId: accountId,
		ruleType:  ruleType,
	}
}

func (r *Rules) Load(ctx context.Context) error {
	if r.ruleType != int64(rule.EventRule) && r.ruleType != int64(rule.BehaviorRule) {
		return microtype.RuleParamFailed
	}

	ac := account.NewAccount(r.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	r.appId = ac.GetQueryAccount().AppID

	do := query.Rule
	mo := do.WithContext(ctx)

	res, err := mo.Where(do.AppID.Eq(r.appId), do.RuleType.Eq(r.ruleType)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return microtype.RuleQueryFailed
	}

	r.res = res

	return nil
}

func (r *Rules) GetResp() *backend.RulesResp {
	rules := make([]*backend.RuleElement, 0, len(r.res))
	for _, r := range r.res {
		if r == nil {
			continue
		}

		rules = append(rules, &backend.RuleElement{
			RuleID:   r.RuleID,
			RuleDesc: r.RuleDesc,
		})
	}

	return &backend.RulesResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		EventRules: rules,
	}
}
