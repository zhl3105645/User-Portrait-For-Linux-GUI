package element

import (
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
)

type UpdateElement struct {
	elementId int64

	req backend.UpdateElementReq
}

func NewUpdateElement(elementId int64, req backend.UpdateElementReq) *UpdateElement {
	return &UpdateElement{
		elementId: elementId,
		req:       req,
	}
}

func (u *UpdateElement) Load(ctx context.Context) error {
	if u.elementId <= 0 {
		return microtype.ParamCheckFailed
	}

	queryMo := struct {
		RuleType int64 `gorm:"column:rule_type;type:int;not null" json:"rule_type"` // 规则类型
	}{}

	err := query.Rule.WithContext(ctx).Select(query.Rule.RuleType).
		Join(query.RuleElement, query.Rule.RuleID.EqCol(query.RuleElement.RuleID)).
		Where(query.RuleElement.RuleElementID.Eq(u.elementId)).
		Scan(&queryMo)
	if err != nil {
		return microtype.ElementQueryFailed
	}

	query.RuleElement.WithContext(ctx).Where(query.RuleElement.RuleElementID)

	mo := model.RuleElement{
		RuleElementValue: "",
	}

	switch queryMo.RuleType {
	case int64(rule.EventRule):
		eventRuleEle := &rule.EventRuleElement{
			EventType:           u.req.EventType,
			MouseClickType:      u.req.MouseClickType,
			MouseClickButton:    u.req.MouseClickButton,
			KeyClickType:        u.req.KeyClickType,
			KeyValue:            u.req.KeyValue,
			ComponentNamePrefix: u.req.ComponentNamePrefix,
		}

		mo.RuleElementValue = rule.GeneEventElement(eventRuleEle)
	case int64(rule.BehaviorRule):
		behaviorRuleEle := &rule.BehaviorRuleElement{
			EventRuleIds: u.req.EventRuleIds,
		}

		mo.RuleElementValue = rule.GeneBehaviorElement(behaviorRuleEle)
	default:

	}

	_, err = query.RuleElement.WithContext(ctx).
		Where(query.RuleElement.RuleElementID.Eq(u.elementId)).
		Updates(mo)
	if err != nil {
		return microtype.ElementUpdateFailed
	}

	return nil
}

func (u *UpdateElement) GetResp() *backend.UpdateElementResp {
	return &backend.UpdateElementResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
