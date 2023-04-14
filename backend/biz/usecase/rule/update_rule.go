package rule

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type UpdateRule struct {
	ruleId  int64
	newDesc string
}

func NewUpdateRule(ruleId int64, newDesc string) *UpdateRule {
	return &UpdateRule{
		ruleId:  ruleId,
		newDesc: newDesc,
	}
}

func (u *UpdateRule) Load(ctx context.Context) error {
	_, err := query.Rule.WithContext(ctx).
		Where(query.Rule.RuleID.Eq(u.ruleId)).
		Update(query.Rule.RuleDesc, u.newDesc)
	if err != nil {
		return microtype.RuleUpdateFailed
	}

	return nil
}

func (u *UpdateRule) GetResp() *backend.UpdateRuleResp {
	return &backend.UpdateRuleResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
