package rule

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DeleteRule struct {
	ruleId int64
}

func NewDeleteRule(ruleId int64) *DeleteRule {
	return &DeleteRule{
		ruleId: ruleId,
	}
}

func (u *DeleteRule) Load(ctx context.Context) error {
	_, err := query.Rule.WithContext(ctx).
		Where(query.Rule.RuleID.Eq(u.ruleId)).
		Delete()
	if err != nil {
		return microtype.RuleDeleteFailed
	}

	return nil
}

func (u *DeleteRule) GetResp() *backend.DeleteRuleResp {
	return &backend.DeleteRuleResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
